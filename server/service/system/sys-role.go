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

// assignMenus 分配菜单
func (s *SysRoleService) assignMenus(tx *gorm.DB, role *system.SysRole, menuIds []uint) ([]system.SysMenu, error) {
	if len(menuIds) == 0 {
		return []system.SysMenu{}, nil
	}
	var menus []system.SysMenu
	// 检查菜单是否存在
	if err := tx.Where("id IN ?", menuIds).Find(&menus).Error; err != nil {
		return nil, errors.New("allMenusNotFound")
	}
	if len(menus) != len(menuIds) {
		return nil, errors.New("menuNotFound")
	}
	// 建立角色菜单关联
	if err := tx.Model(role).Association("Menus").Replace(menus); err != nil {
		return nil, errors.New("menuAssignFailed")
	}

	// 获取菜单下绑定的api
	if err := tx.Model(&system.SysMenu{}).
		Preload("Apis").
		Where("id IN ?", lo.Map(menus, func(item system.SysMenu, index int) uint {
			return item.Id
		})).
		Find(&menus).
		Error; err != nil {
		return nil, errors.New("menuAssignFailed")
	}
	// 将菜单下的APIs合并
	apis := lo.FlatMap(menus, func(menu system.SysMenu, _ int) []system.SysApi {
		return menu.Apis
	})

	if len(apis) == 0 { // 说明没有分配API
		return []system.SysMenu{}, nil
	}

	enforcer := global.Enforcer

	if _, err := enforcer.RemoveFilteredPolicy(0, role.Code); err != nil {
		return nil, errors.New("menuAssignFailed")
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
	if err != nil {
		return nil, errors.New("menuAssignFailed")
	}
	if !success {
		return nil, errors.New("menuAssignFailed")
	}
	return menus, nil
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
func (s *SysRoleService) checkCodeDuplicate(code string) error {
	var count int64
	err := global.DB.Model(&system.SysRole{}).
		Where("code = ?", code).
		Count(&count).
		Error
	if err != nil || count > 0 {
		return errors.New("roleCodeDuplicated")
	}
	return nil
}

// Create 创建角色
func (s *SysRoleService) Create(req request.CreateSysRoleReq) (*response.SysRoleResp, error) {
	db := global.DB
	if err := s.checkCodeDuplicate(req.Code); err != nil {
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

// Delete 删除角色
func (s *SysRoleService) Delete(id uint) error {
	db := global.DB
	// 先检查角色是否存在
	existRole, err := s.checkRoleExist(db, id)
	if err != nil {
		return err
	}

	err = db.
		Delete(existRole, id).
		Error
	if err != nil {
		return errors.New("deleteRoleFailed")
	}

	enforcer := global.Enforcer
	if _, err := enforcer.RemoveFilteredPolicy(0, existRole.Code); err != nil {
		return errors.New("deleteRoleFailed")
	}

	return nil
}

// Update 更新角色
func (s *SysRoleService) Update(req request.UpdateSysRoleReq) error {
	db := global.DB
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

	if _, err := s.assignMenus(tx, role, req.MenuIds); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("updateRoleFailed")
	}
	return nil
}

// Get 查询单个角色
func (s *SysRoleService) Get(id uint) (*response.SysRoleResp, error) {
	db := global.DB

	_, err := s.checkRoleExist(db, id)
	if err != nil {
		return nil, err
	}

	var role system.SysRole
	err = db.
		Preload("Menus").
		First(&role, id).
		Error
	if err != nil {
		return nil, errors.New("getRoleFailed")
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
func (s *SysRoleService) GetList(req request.QuerySysRoleReq) (common.PageResult[response.SysRoleResp], error) {
	var roles []system.SysRole
	var total int64
	db := global.DB.Model(&system.SysRole{})

	// 添加模糊查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
	}

	db.Count(&total)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	db = db.
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Order("id ASC")
	// 添加预加载菜单数据
	err := db.Preload("Menus").Find(&roles).Error
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
func (s *SysRoleService) GetAll() ([]response.SysRoleShortResp, error) {
	var roles []system.SysRole
	db := global.DB.Model(&system.SysRole{})
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
