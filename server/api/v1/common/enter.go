package common

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	AuthApi
	UserApi
	OssApi
}

var (
	authService         = service.ServiceGroupApp.CommonServiceGroup.AuthService      // 认证服务
	userService         = service.ServiceGroupApp.CommonServiceGroup.UserService      // 用户服务
	ossService          = service.ServiceGroupApp.CommonServiceGroup.OssService       // OSS服务
	tenantPublicService = service.ServiceGroupApp.SystemServiceGroup.SysTenantService // 租户服务（登录页只读列表）
)
