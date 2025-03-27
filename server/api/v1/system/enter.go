package system

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	SysUserApi
}

var (
	userService = service.ServiceGroupApp.SystemServiceGroup.SysUserService // 系统用户服务
)
