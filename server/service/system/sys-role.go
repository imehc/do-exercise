package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type SysRoleService struct{}

// assignMenus 分配菜单。
// 空 menuIds、以及菜单未绑定任何 API 这两种情况，语义都是「撤销该角色的全部策略」。
// 早期实现在这两处提前 return，跳过了 RemoveFilteredPolicy，
// 导致解除权限的请求返回成功而 Casbin 策略原样保留。
// 本函数只负责在事务内落库角色-菜单关联表；Casbin p 策略的同步
// 放在事务提交之后由 syncRolePolicy 统一处理。
func (s *SysRoleService) assignMenus(tx *gorm.DB, role *system.SysRole, menuIds []uint) ([]system.SysMenu, error) {
	var menus []system.SysMenu
	if len(menuIds) > 0 {
		// 一次查出菜单及其绑定的 API（此前分两次查询同一批行，第二次只为补 Apis）
		if err := tx.Preload("Apis").Where("id IN ?", menuIds).Find(&menus).Error; err != nil {
			return nil, errors.New("allMenusNotFound")
		}
		if len(menus) != len(menuIds) {
			return nil, errors.New("menuNotFound")
		}
	}

	// 建立/清空角色菜单关联
	if len(menus) == 0 {
		if err := tx.Model(role).Association("Menus").Clear(); err != nil {
			return nil, errors.New("menuAssignFailed")
		}
	} else if err := tx.Model(role).Association("Menus").Replace(menus); err != nil {
		return nil, errors.New("menuAssignFailed")
	}

	return menus, nil
}

// syncRolePolicy 在事务提交成功后同步 Casbin p 策略。
// global.Enforcer 走独立 adapter，无法加入业务事务，必须等 Commit() 成功后再写，
// 否则提交失败时角色行回滚、策略却已生效（幽灵授权）。
// 此步失败则角色表现为无权限（fail closed），并尽量回退半成品规则。
func (s *SysRoleService) syncRolePolicy(role *system.SysRole, menus []system.SysMenu) error {
	enforcer := global.Enforcer
	// 无条件清空旧策略，再按新集合重建
	if _, err := enforcer.RemoveFilteredPolicy(0, role.Code); err != nil {
		return errors.New("menuAssignFailed")
	}

	// 将菜单下的APIs合并
	apis := lo.FlatMap(menus, func(menu system.SysMenu, _ int) []system.SysApi {
		return menu.Apis
	})
	if len(apis) == 0 {
		// 菜单已分配但未绑定 API：策略为空，菜单列表仍需如实返回
		return nil
	}

	// 使用casbin批量添加策略
	policies := lo.Map(apis, func(item system.SysApi, index int) []string {
		return []string{
			role.Code,
			item.Path,
			item.Method,
		}
	})

	// 添加策略并检查结果
	success, err := enforcer.AddPolicies(policies)
	if err != nil || !success {
		// 补偿：清掉可能已加入的部分策略，回退到无权限态
		_, _ = enforcer.RemoveFilteredPolicy(0, role.Code)
		return errors.New("menuAssignFailed")
	}
	return nil
}

// checkRoleExist 检查角色是否存在
func (s *SysRoleService) checkRoleExist(db *gorm.DB, roleId uint) (*system.SysRole, error) {
	var role system.SysRole
	result := db.
		Unscoped().
		First(&role, roleId)
	if result.Error != nil {
		return nil, errors.New("allRolesNotFound")
	}

	if !role.DeletedAt.Time.IsZero() {
		return nil, errors.New("roleDeleted")
	}

	return &role, nil
}

// checkCodeDuplicate 检查角色编码是否重复
func (s *SysRoleService) checkCodeDuplicate(db *gorm.DB, code string) error {
	var count int64
	err := db.Model(&system.SysRole{}).
		Where("code = ?", code).
		Count(&count).
		Error
	if err != nil || count > 0 {
		return errors.New("roleCodeDuplicated")
	}
	return nil
}

// Create 创建角色
func (s *SysRoleService) Create(db *gorm.DB, req request.CreateSysRoleReq) (*response.SysRoleResp, error) {
	if err := s.checkCodeDuplicate(db, req.Code); err != nil {
		return nil, err
	}

	role := &system.SysRole{
		Name: req.Name,
		Code: req.Code,
	}

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createRoleFailed")
	}

	menus, err := s.assignMenus(tx, role, req.MenuIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("createRoleFailed")
	}

	// Casbin 策略在事务提交成功之后同步
	if err := s.syncRolePolicy(role, menus); err != nil {
		return nil, err
	}

	return &response.SysRoleResp{
		Id:        role.Id,
		Name:      role.Name,
		Code:      role.Code,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		Menus: lo.Map(menus, func(item system.SysMenu, index int) response.SysMenuShortResp {
			return response.SysMenuShortResp{
				Id:   item.Id,
				Name: item.Name,
			}
		}),
	}, nil
}

// Delete 删除角色。
// 删除记录、解除用户关联、清理 Casbin 策略三步必须一致：
// 早期实现无事务且不清理 g 规则，一旦策略删除失败，角色虽被软删，
// 但持有该 role code 的用户仍保留全部权限。
func (s *SysRoleService) Delete(db *gorm.DB, id uint) error {
	// 先检查角色是否存在
	existRole, err := s.checkRoleExist(db, id)
	if err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id = ?", id).Delete(existRole).Error; err != nil {
		tx.Rollback()
		return errors.New("deleteRoleFailed")
	}

	// 解除角色与菜单的关联，避免留下悬空的中间表数据
	if err := tx.Model(existRole).Association("Menus").Clear(); err != nil {
		tx.Rollback()
		return errors.New("deleteRoleFailed")
	}
	// SysRole 未声明 Users 反向关联，用户侧的中间表需直接清理。
	// 列名 sys_role_id 由 GORM 命名策略生成（join table: sys_user_role）。
	if err := tx.Table("sys_user_role").
		Where("sys_role_id = ?", existRole.Id).
		Delete(nil).Error; err != nil {
		tx.Rollback()
		return errors.New("deleteRoleFailed")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("deleteRoleFailed")
	}

	// Casbin 走独立适配器，不在上面的事务内，故放到提交成功之后
	enforcer := global.Enforcer
	// p 策略：该角色拥有的权限
	if _, err := enforcer.RemoveFilteredPolicy(0, existRole.Code); err != nil {
		return errors.New("deleteRoleFailed")
	}
	// g 规则：仍指向该角色的用户绑定，不清理会永久残留
	if _, err := enforcer.RemoveFilteredGroupingPolicy(1, existRole.Code); err != nil {
		return errors.New("deleteRoleFailed")
	}

	return nil
}

// Update 更新角色
func (s *SysRoleService) Update(db *gorm.DB, req request.UpdateSysRoleReq) error {
	role, err := s.checkRoleExist(db, req.Id)
	if err != nil {
		return err
	}
	role.Name = req.Name

	// 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.
		Model(role).
		Omit("id", "created_at", "created_by").
		Updates(&role).
		Error; err != nil {
		tx.Rollback()
		return errors.New("updateRoleFailed")
	}

	menus, err := s.assignMenus(tx, role, req.MenuIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("updateRoleFailed")
	}

	// 事务提交成功后在 Casbin 同步策略
	if err := s.syncRolePolicy(role, menus); err != nil {
		return err
	}
	return nil
}

// Get 查询单个角色
func (s *SysRoleService) Get(db *gorm.DB, id uint) (*response.SysRoleResp, error) {
	// 一次查询带出菜单，不再先查存在性再重查同一行
	var role system.SysRole
	result := db.
		Unscoped().
		Preload("Menus").
		First(&role, id)
	if result.Error != nil {
		return nil, errors.New("allRolesNotFound")
	}
	if !role.DeletedAt.Time.IsZero() {
		return nil, errors.New("roleDeleted")
	}

	menus := make([]response.SysMenuShortResp, len(role.Menus))
	for i, menu := range role.Menus {
		menus[i] = response.SysMenuShortResp{
			Id:   menu.Id,
			Name: menu.Name,
		}
	}
	return &response.SysRoleResp{
		Id:        role.Id,
		Name:      role.Name,
		Code:      role.Code,
		CreatedAt: role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
		Menus:     menus,
	}, nil
}

// GetList 查询角色列表
func (s *SysRoleService) GetList(db *gorm.DB, req request.QuerySysRoleReq) (common.PageResult[response.SysRoleResp], error) {
	var roles []system.SysRole
	var total int64

	// Count 用独立 builder，避免污染后续 Find 的状态
	countDB := db.Model(&system.SysRole{})

	// 添加模糊查询条件
	if req.Name != "" {
		countDB = countDB.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		countDB = countDB.Where("code LIKE ?", "%"+req.Code+"%")
	}

	if err := countDB.Count(&total).Error; err != nil {
		return common.PageResult[response.SysRoleResp]{}, errors.New("getRoleListFailed")
	}
	req.Normalize()
	err := db.Model(&system.SysRole{}).
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Order("id ASC").
		Preload("Menus").
		Find(&roles).Error
	if err != nil {
		return common.PageResult[response.SysRoleResp]{}, errors.New("getRoleListFailed")
	}
	data := make([]response.SysRoleResp, len(roles))
	for i, role := range roles {
		// 转换菜单数据
		menus := make([]response.SysMenuShortResp, len(role.Menus))
		for j, menu := range role.Menus {
			menus[j] = response.SysMenuShortResp{
				Id:   menu.Id,
				Name: menu.Name,
			}
		}

		data[i] = response.SysRoleResp{
			Id:        role.Id,
			Name:      role.Name,
			Code:      role.Code,
			CreatedAt: role.CreatedAt,
			CreatedBy: role.CreatedBy,
			UpdatedAt: role.UpdatedAt,
			UpdatedBy: role.UpdatedBy,
			Menus:     menus,
		}
	}
	result := common.PageResult[response.SysRoleResp]{
		Data: data,
		Meta: common.PageMeta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	}
	return result, nil
}

// GetAll 获取所有角色
func (s *SysRoleService) GetAll(db *gorm.DB) ([]response.SysRoleShortResp, error) {
	var roles []system.SysRole
	db = db.Model(&system.SysRole{})
	err := db.
		Order("id ASC").
		Find(&roles).
		Error
	if err != nil {
		return nil, errors.New("getRoleFailed")
	}

	return lo.Map(roles, func(role system.SysRole, _ int) response.SysRoleShortResp {
		return response.SysRoleShortResp{
			Id:   role.Id,
			Code: role.Code,
			Name: role.Name,
		}
	}), nil
}
