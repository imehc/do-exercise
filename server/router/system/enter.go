package system

import api "github.com/imehc/do-exercise/server/api/v1"

type RouterGroup struct {
	SysUserRouter
	SysMenuRouter
	SysRoleRouter
	SysApiRouter
	SysOperationLogRouter
	SysTokenRouter
	SysInfoRouter
	SysJobRouter
	SysTenantRouter
}

var (
	sysUserApi         = api.ApiGroupApp.SystemApiGroup.SysUserApi
	sysMenuApi         = api.ApiGroupApp.SystemApiGroup.SysMenuApi
	sysRoleApi         = api.ApiGroupApp.SystemApiGroup.SysRoleApi
	sysApiApi          = api.ApiGroupApp.SystemApiGroup.SysApiApi
	sysOperationLogApi = api.ApiGroupApp.SystemApiGroup.SysOperationLogApi
	sysTokenApi        = api.ApiGroupApp.SystemApiGroup.SysTokenApi
	sysInfoApi         = api.ApiGroupApp.SystemApiGroup.SysInfoApi
	sysJobApi          = api.ApiGroupApp.SystemApiGroup.SysJobApi
	tenantApi          = api.ApiGroupApp.SystemApiGroup.SysTenantApi
)
