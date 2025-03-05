package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/pkg/utils/scope"
	"gorm.io/gorm"
)

type RoleService struct{}

// 创建角色
func (r RoleService) Create(request request.CreateRoleRequest, createBy uint) (err error) {
	db := global.DB

	role := system.Role{
		Name: request.Name,
		Key:  request.Key,
		ControlWrapper: model.ControlWrapper{
			CreateBy: createBy,
		},
	}

	if request.Sort != 0 {
		role.Sort = request.Sort
	}
	if request.Status != 0 {
		role.Status = request.Status
	}
	if request.Remark != "" {
		role.Remark = request.Remark
	}

	err = db.Create(&role).Error
	return
}

// 删除角色
func (r RoleService) Delete(param request.RoleParam, deleteBy uint) (err error) {
	db := global.DB

	var role system.Role
	result := db.
		Unscoped().
		First(&role, param.RoleId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("角色不存在")
		}
		return result.Error
	}

	if !role.DeleteAt.Time.IsZero() {
		return errors.New("角色已删除")
	}

	// 检查是否有权限删除
	if role.CreateBy != deleteBy {
		// 检查是否是超级管理员
		var currentUserRole system.Role
		if err := db.First(&currentUserRole, deleteBy).Error; err != nil {
			return err
		}
		if !currentUserRole.IsAdmin {
			return errors.New("无权删除其他用户创建的数据")
		}
	}

	var user system.User
	result = db.
		First(&user, "role_id = ?", param.RoleId)
	if result.Error == nil && user.UserId != 0 {
		return errors.New("该角色已被使用,无法删除")
	}

	db.
		Model(system.Role{}).
		Where("role_id = ?", param.RoleId).
		Update("deleted_by", deleteBy).
		Delete(&role)
	return nil
}

// 更新角色
func (r RoleService) Update(param request.RoleParam, request request.UpdateRoleRequest, updateBy uint) (err error) {
	db := global.DB

	var role system.Role
	result := db.
		First(&role, param.RoleId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("角色不存在")
		}
		return result.Error
	}

	role.Sort = request.Sort
	role.Status = request.Status
	role.Remark = request.Remark
	role.ControlWrapper = model.ControlWrapper{
		UpdateBy: updateBy,
	}

	db.
		Model(system.Role{}).
		Where("role_id = ?", param.RoleId).
		Updates(&role).
		Omit("role_id", "created_at")

	return nil
}

// 查询角色
func (r RoleService) Find(param request.RoleParam) (response sysRes.RoleItem, err error) {
	db := global.DB

	var role system.Role
	result := db.
		Preload("Menus").
		Preload("Apis").
		First(&role, param.RoleId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return response, errors.New("角色不存在")
		}
		return response, result.Error
	}

	response.ID = role.RoleId
	response.Name = role.Name
	response.Key = role.Key
	response.IsAdmin = role.IsAdmin
	response.DataScope = role.DataScope
	response.ControlWrapper = role.ControlWrapper
	response.Sort = role.Sort
	response.Status = role.Status
	response.Remark = role.Remark

	// 加载菜单信息
	response.Menus = make([]sysRes.MenuItem, len(role.Menus))
	for i, menu := range role.Menus {
		response.Menus[i] = sysRes.MenuItem{
			IDWrapper: common.IDWrapper{
				ID: menu.MenuId,
			},
			ControlWrapper: menu.ControlWrapper,
			MenuRequest: request.MenuRequest{
				ParentId:   &menu.ParentId,
				Name:       menu.Name,
				Icon:       menu.Icon,
				Type:       menu.Type,
				Action:     menu.Action,
				Title:      menu.Title,
				Path:       menu.Path,
				Permission: menu.Permission,
				SortWrapper: model.SortWrapper{
					Sort: menu.Sort,
				},
			},
			Route: menu.Route,
		}
	}

	// 加载API信息
	response.Apis = make([]sysRes.ApiItem, len(role.Apis))
	for i, api := range role.Apis {
		response.Apis[i] = sysRes.ApiItem{
			IDWrapper: common.IDWrapper{
				ID: api.ApiId,
			},
			ControlWrapper: api.ControlWrapper,
			ApiRequest: request.ApiRequest{
				Handle: api.Handle,
				Title:  api.Title,
				Path:   api.Path,
				Type:   api.Type,
				Action: api.Action,
			},
		}
	}

	return
}

// 查询角色列表
func (r RoleService) FindList(query request.RoleQueryParams, s common.ScopeData) (response sysRes.RoleResponse, err error) {
	db := global.DB
	// 应用数据权限过滤
	db = scope.GetDataScope(db, &s, "sys_role")

	var total int64
	var originRoles []system.Role
	db = db.Model(&system.Role{}).
		Order("sort Desc").
		Order("role_id ASC").
		Preload("Menus").
		Preload("Apis").
		Where("name LIKE ?", "%"+query.Name+"%")

	err = db.Count(&total).
		Find(&originRoles).
		Error
	if err != nil {
		return response, err
	}

	response.Meta.Page = query.Page
	response.Meta.PageSize = query.PageSize
	response.Meta.Total = total

	response.Data = make([]sysRes.RoleItem, len(originRoles))
	for i, role := range originRoles {
		response.Data[i].ID = role.RoleId
		response.Data[i].ControlWrapper = role.ControlWrapper
		response.Data[i].Name = role.Name
		response.Data[i].Key = role.Key
		response.Data[i].IsAdmin = role.IsAdmin
		response.Data[i].DataScope = role.DataScope
		response.Data[i].Sort = role.Sort
		response.Data[i].Status = role.Status
		response.Data[i].Remark = role.Remark

		// 加载菜单信息
		response.Data[i].Menus = make([]sysRes.MenuItem, len(role.Menus))
		for j, menu := range role.Menus {
			response.Data[i].Menus[j] = sysRes.MenuItem{
				IDWrapper: common.IDWrapper{
					ID: menu.MenuId,
				},
				ControlWrapper: menu.ControlWrapper,
				MenuRequest: request.MenuRequest{
					ParentId:   &menu.ParentId,
					Name:       menu.Name,
					Icon:       menu.Icon,
					Type:       menu.Type,
					Action:     menu.Action,
					Title:      menu.Title,
					Path:       menu.Path,
					Permission: menu.Permission,
					SortWrapper: model.SortWrapper{
						Sort: menu.Sort,
					},
				},
				Route: menu.Route,
			}
		}

		// 加载API信息
		response.Data[i].Apis = make([]sysRes.ApiItem, len(role.Apis))
		for j, api := range role.Apis {
			response.Data[i].Apis[j] = sysRes.ApiItem{
				IDWrapper: common.IDWrapper{
					ID: api.ApiId,
				},
				ControlWrapper: api.ControlWrapper,
				ApiRequest: request.ApiRequest{
					Handle: api.Handle,
					Title:  api.Title,
					Path:   api.Path,
					Type:   api.Type,
					Action: api.Action,
				},
			}
		}
	}

	return
}

// 更新角色数据权限
func (r RoleService) UpdateDataScope(param request.RoleParam, request request.UpdateRoleDataScope, updateBy uint) (err error) {
	db := global.DB

	var role system.Role
	result := db.
		First(&role, param.RoleId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("角色不存在")
		}
		return result.Error
	}

	// 更新角色的数据范围
	role.DataScope = request.DataScope
	role.ControlWrapper = model.ControlWrapper{
		UpdateBy: updateBy,
	}

	// 如果是自定义数据权限，需要处理角色与部门的关联
	if request.DataScope == 2 { // 自定数据权限
		// 检查部门是否存在
		if len(request.DeptIds) > 0 {
			var depts []system.Dept
			if err = db.Where("dept_id IN ?", request.DeptIds).Find(&depts).Error; err != nil {
				return err
			}
			if len(depts) != len(request.DeptIds) {
				return errors.New("部分部门不存在")
			}

			// 更新角色与部门的关联关系
			if err = db.Model(&role).Association("Depts").Replace(depts); err != nil {
				return err
			}
		} else {
			// 如果没有指定部门，清除所有关联关系
			if err = db.Model(&role).Association("Depts").Clear(); err != nil {
				return err
			}
		}
	} else {
		// 非自定义数据权限时，清除所有部门关联
		if err = db.Model(&role).Association("Depts").Clear(); err != nil {
			return err
		}
	}

	// 更新角色基本信息
	db.
		Model(system.Role{}).
		Where("role_id = ?", param.RoleId).
		Updates(&role).
		Omit("role_id", "created_at")

	return nil
}

// 更新角色菜单权限
func (r RoleService) UpdateMenuScope(param request.RoleParam, request request.UpdateMenuDataScope, updateBy uint) (err error) {
	db := global.DB

	var role system.Role
	result := db.
		First(&role, param.RoleId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("角色不存在")
		}
		return result.Error
	}

	// 检查菜单是否存在
	if len(request.MenuIds) > 0 {
		var menus []system.Menu
		if err = db.Where("menu_id IN ?", request.MenuIds).Find(&menus).Error; err != nil {
			return err
		}
		if len(menus) != len(request.MenuIds) {
			return errors.New("部分菜单不存在")
		}

		// 更新角色与菜单的关联关系
		if err = db.Model(&role).Association("Menus").Replace(menus); err != nil {
			return err
		}
	} else {
		// 如果没有指定菜单，清除所有关联关系
		if err = db.Model(&role).Association("Menus").Clear(); err != nil {
			return err
		}
	}

	// 更新角色基本信息
	role.ControlWrapper = model.ControlWrapper{
		UpdateBy: updateBy,
	}
	db.
		Model(system.Role{}).
		Where("role_id = ?", param.RoleId).
		Updates(&role).
		Omit("role_id", "created_at")

	return nil
}

// 更新角色api权限
func (r RoleService) UpdateApiScope(param request.RoleParam, request request.UpdateApiDataScope, updateBy uint) (err error) {
	db := global.DB

	var role system.Role
	result := db.
		First(&role, param.RoleId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("角色不存在")
		}
		return result.Error
	}

	// 检查API是否存在
	if len(request.ApiIds) > 0 {
		var apis []system.Api
		if err = db.Where("api_id IN ?", request.ApiIds).Find(&apis).Error; err != nil {
			return err
		}
		if len(apis) != len(request.ApiIds) {
			return errors.New("部分API不存在")
		}

		// 更新角色与API的关联关系
		if err = db.Model(&role).Association("Apis").Replace(apis); err != nil {
			return err
		}
	} else {
		// 如果没有指定API，清除所有关联关系
		if err = db.Model(&role).Association("Apis").Clear(); err != nil {
			return err
		}
	}

	// 更新角色基本信息
	role.ControlWrapper = model.ControlWrapper{
		UpdateBy: updateBy,
	}
	db.
		Model(system.Role{}).
		Where("role_id = ?", param.RoleId).
		Updates(&role).
		Omit("role_id", "created_at")

	return nil
}
