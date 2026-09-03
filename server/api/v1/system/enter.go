package system

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	SysUserApi
	SysMenuApi
	SysRoleApi
	SysApiApi
	SysOperationLogApi
	SysTokenApi
	SysInfoApi
	SysJobApi
	SysTenantApi
}

var (
	userService            = service.ServiceGroupApp.SystemServiceGroup.SysUserService         // 系统用户服务
	menuService            = service.ServiceGroupApp.SystemServiceGroup.SysMenuService         // 系统菜单服务
	roleService            = service.ServiceGroupApp.SystemServiceGroup.SysRoleService         // 系统角色服务
	apiService             = service.ServiceGroupApp.SystemServiceGroup.SysApiService          // 系统api服务
	sysOperationLogService = service.ServiceGroupApp.SystemServiceGroup.SysOperationLogService // 系统操作日志服务
	sysTokenService        = service.ServiceGroupApp.SystemServiceGroup.SysTokenService        // 令牌服务
	sysInfoService         = service.ServiceGroupApp.SystemServiceGroup.SysInfoService         // 系统信息服务
	jobService             = service.ServiceGroupApp.SystemServiceGroup.SysJobService          // 定时任务服务
	tenantService          = service.ServiceGroupApp.SystemServiceGroup.SysTenantService       // 租户服务
)
