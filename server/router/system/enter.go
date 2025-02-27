package system

import api "github.com/imehc/do-exercise/server/api/v1"

type RouterGroup struct {
	AuthRouter
}

var (
	authApi    = api.ApiGroupApp.SystemApiGroup.AuthApi
	captchaApi = api.ApiGroupApp.SystemApiGroup.CaptchaApi
	deptApi    = api.ApiGroupApp.SystemApiGroup.DeptApi
	dictApi    = api.ApiGroupApp.SystemApiGroup.DictApi
	menuApi    = api.ApiGroupApp.SystemApiGroup.MenuApi
	postApi    = api.ApiGroupApp.SystemApiGroup.PostApi
	roleApi    = api.ApiGroupApp.SystemApiGroup.RoleApi
	userApi    = api.ApiGroupApp.SystemApiGroup.UserApi
	apiApi     = api.ApiGroupApp.SystemApiGroup.ApiApi
)
