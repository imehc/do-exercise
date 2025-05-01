package common

import api "github.com/imehc/do-exercise/server/api/v1"

type RouterGroup struct {
	AuthRouter
	UserRouter
	OssRouter
}

var (
	authApi = api.ApiGroupApp.CommonApiGroup.AuthApi
	userApi = api.ApiGroupApp.CommonApiGroup.UserApi
	ossApi  = api.ApiGroupApp.CommonApiGroup.OssApi
)
