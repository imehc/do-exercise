package system

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	SysUserApi
	SysMenuApi
	SysRoleApi
	SysApiApi
	SysOperationLogApi
}

var (
	userService            = service.ServiceGroupApp.SystemServiceGroup.SysUserService         // 系统用户服务
	menuService            = service.ServiceGroupApp.SystemServiceGroup.SysMenuService         // 系统菜单服务
	roleService            = service.ServiceGroupApp.SystemServiceGroup.SysRoleService         // 系统角色服务
	apiService             = service.ServiceGroupApp.SystemServiceGroup.SysApiService          // 系统api服务
	sysOperationLogService = service.ServiceGroupApp.SystemServiceGroup.SysOperationLogService // 系统操作日志服务
)
