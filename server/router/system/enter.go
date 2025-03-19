package system

import api "github.com/imehc/do-exercise/server/api/v1"

type RouterGroup struct {
	SysUserRouter
}

var (
	sysUserApi = api.ApiGroupApp.SystemApiGroup.SysUserApi
)
