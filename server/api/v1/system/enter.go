package system

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	AuthApi
	CaptchaApi
}

var (
	authService = service.ServiceGroupApp.SystemServiceGroup.AuthService // 登录服务
)
