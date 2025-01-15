package system

import api "github.com/imehc/do-exercise/server/api/v1"

type RouterGroup struct {
	AuthRouter
}

var (
	authApi    = api.ApiGroupApp.SystemApiGroup.AuthApi
	captchaApi = api.ApiGroupApp.SystemApiGroup.CaptchaApi
)
