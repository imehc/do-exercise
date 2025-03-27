package common

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	AuthApi
}

var (
	userService = service.ServiceGroupApp.SystemServiceGroup.SysUserService // 系统用户服务
)
