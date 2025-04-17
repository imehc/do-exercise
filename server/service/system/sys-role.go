package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
)

type SysRoleService struct{}

// Create 创建角色
func (s *SysRoleService) Create(req request.CreateSysRoleReq) (*response.SysRoleResp, error) {
	role := &system.SysRole{
		Name: req.Name,
		Code: req.Code,
	}
	err := global.DB.Create(role).Error
	if err != nil {
		return nil, err
	}
	return &response.SysRoleResp{
		Id:        role.Id,
		Name:      role.Name,
		Code:      role.Code,
		CreatedAt: role.CreatedAt,
		CreatedBy: role.CreatedBy,
		UpdatedAt: role.UpdatedAt,
		UpdatedBy: role.UpdatedBy,
	}, nil
}

// Delete 删除角色
func (s *SysRoleService) Delete(id uint) error {
	// 先检查角色是否存在
	existRole := &system.SysRole{}
	err := global.DB.Where("id = ?", id).
		First(existRole).
		Error
	if err != nil {
		return err
	}
	return global.DB.
		Delete(&system.SysRole{}, id).
		Error
}

// Update 更新角色
func (s *SysRoleService) Update(req request.UpdateSysRoleReq) error {
	// 先检查角色是否存在
	var role system.SysRole
	result := global.DB.
		First(&role, req.Id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("菜单不存在")
		}
		return result.Error
	}
	role.Name = req.Name

	return global.DB.
		Model(system.SysRole{}).
		Where("id = ?", req.Id).
		Updates(&role).
		Omit("id", "created_at", "created_by").
		Error
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
