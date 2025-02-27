package system

import (
	"github.com/imehc/do-exercise/server/model/system/request"
	sysRes "github.com/imehc/do-exercise/server/model/system/response"
)

type MenuService struct{}

// 创建菜单
func (m *MenuService) CreateMenu(request request.MenuRequest, createdBy uint) (err error) {
	return
}

// 删除菜单
func (m *MenuService) DeleteMenu(param request.MenuParam, deletedBy uint) (err error) {
	return
}

// 更新菜单
func (m *MenuService) UpdateMenu(param request.MenuParam, request request.MenuRequest, updatedBy uint) (err error) {
	return
}

// 获取菜单
func (m *MenuService) GetMenu(param request.MenuParam) (response []sysRes.MenuItem, err error) {
	return
}

// 获取菜单列表
func (m *MenuService) GetMenuList(query request.MenuQueryParams) (response []sysRes.MenuResponse, err error) {
	return
}

// 获取菜单树
func (m *MenuService) GetMenuTree() (response []sysRes.MenuTree, err error) {
	return
}
