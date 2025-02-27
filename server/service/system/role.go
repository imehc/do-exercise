package system

import (
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
)

type RoleService struct{}

// 创建角色
func (r RoleService) Create(request request.CreateRoleRequest, createdBy uint) (err error) {
	return
}

// 删除角色
func (r RoleService) Delete(param request.RoleParam, deletedBy uint) (err error) {
	return
}

// 更新角色
func (r RoleService) Update(param request.RoleParam, request request.UpdateRoleRequest, updatedBy uint) (err error) {
	return
}

// 查询角色
func (r RoleService) Find(param request.RoleParam) (response sysRes.RoleItem, err error) {
	return
}

// 查询角色列表
func (r RoleService) FindList(query request.RoleQueryParams) (response []sysRes.RoleResponse, err error) {
	return
}

// 更新角色数据权限
func (r RoleService) UpdateDataScope(param request.RoleParam, request request.UpdateRoleDataScope, updatedBy uint) (err error) {
	return
}

// 更新角色菜单权限
func (r RoleService) UpdateMenuScope(param request.RoleParam, request request.UpdateMenuDataScope, updatedBy uint) (err error) {
	return
}

// 更新角色api权限
func (r RoleService) UpdateApiScope(param request.RoleParam, request request.UpdateApiDataScope, updatedBy uint) (err error) {
	return
}
