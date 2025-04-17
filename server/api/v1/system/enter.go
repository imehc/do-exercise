package system

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	SysUserApi
	SysMenuApi
}

var (
	userService = service.ServiceGroupApp.SystemServiceGroup.SysUserService // 系统用户服务
	menuService = service.ServiceGroupApp.SystemServiceGroup.SysMenuService // 系统菜单服务
)
