package common

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	AuthApi
	UserApi
}

var (
	authService = service.ServiceGroupApp.CommonServiceGroup.AuthService // 认证服务
	userService = service.ServiceGroupApp.CommonServiceGroup.UserService // 用户服务
)
