package common

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	AuthApi
}

var (
	authService = service.ServiceGroupApp.CommonServiceGroup.AuthService // 认证服务
)
