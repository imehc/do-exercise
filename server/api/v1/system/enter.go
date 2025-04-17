package system

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	SysUserApi
	SysMenuApi
	SysRoleApi
}

var (
	userService = service.ServiceGroupApp.SystemServiceGroup.SysUserService // 系统用户服务
	menuService = service.ServiceGroupApp.SystemServiceGroup.SysMenuService // 系统菜单服务
	roleService = service.ServiceGroupApp.SystemServiceGroup.SysRoleService // 角色菜单服务
)
