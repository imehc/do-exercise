package system

import (
	"errors"
	"fmt"

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
		return nil, err
	}
	if len(menus) != len(menuIds) {
		return nil, errors.New("菜单不存在")
	}
	// 建立角色菜单关联
	if err := tx.Model(role).Association("Menus").Replace(menus); err != nil {
		return nil, err
	}

	// 获取菜单下绑定的api
	if err := tx.Model(&system.SysMenu{}).
		Preload("Apis").
		Where("id IN ?", lo.Map(menus, func(item system.SysMenu, index int) uint {
			return item.Id
		})).
		Find(&menus).
		Error; err != nil {
		return nil, err
	}
	// 将菜单下的APIs合并
	apis := lo.FlatMap(menus, func(menu system.SysMenu, _ int) []system.SysApi {
		return menu.Apis
	})

	if len(apis) == 0 {
		return menus, nil
	}

	enforcer := global.Enforcer

	if _, err := enforcer.RemoveFilteredPolicy(0, role.Code); err != nil {
		return nil, errors.New(fmt.Sprintf("清除策略失败: %v", err))
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
		return nil, errors.New(fmt.Sprintf("添加策略失败: %v", err))
	}
	if !success {
		return nil, errors.New("部分策略添加失败")
	}
	return menus, nil
}

// Create 创建角色
func (s *SysRoleService) Create(req request.CreateSysRoleReq) (*response.SysRoleResp, error) {
	role := &system.SysRole{
		Name: req.Name,
		Code: req.Code,
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(role).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	menus, err := s.assignMenus(tx, role, req.MenuIds)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
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
		})}, nil
}

// Delete 删除角色
func (s *SysRoleService) Delete(id uint) error {
	db := global.DB
	// 先检查角色是否存在
	existRole := &system.SysRole{}
	err := db.Where("id = ?", id).
		First(existRole).
		Error
	if err != nil {
		return err
	}
	err = db.
		Delete(&system.SysRole{}, id).
		Error
	if err != nil {
		return err
	}

	enforcer := global.Enforcer
	if _, err := enforcer.RemoveFilteredPolicy(0, existRole.Code); err != nil {
		return errors.New(fmt.Sprintf("清除策略失败: %v", err))
	}

	return nil
}

// Update 更新角色
func (s *SysRoleService) Update(req request.UpdateSysRoleReq) error {
	db := global.DB
	// 先检查角色是否存在
	var role system.SysRole
	result := db.
		First(&role, req.Id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("菜单不存在")
		}
		return result.Error
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
		Model(system.SysRole{}).
		Where("id = ?", req.Id).
		Updates(&role).
		Omit("id", "created_at", "created_by").
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if _, err := s.assignMenus(tx, &role, req.MenuIds); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Get 查询单个角色
func (s *SysRoleService) Get(id uint) (*response.SysRoleResp, error) {
	role := &system.SysRole{}
	err := global.DB.Where("id =?", id).
		Preload("Menus").
		First(role).
		Error
	if err != nil {
		return nil, err
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
	db = db.Scopes(util.Paginate(req.PageSize, req.Page))
	// 添加预加载菜单数据
	err := db.Preload("Menus").Find(&roles).Error
	if err != nil {
		return common.PageResult[response.SysRoleResp]{}, err
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
