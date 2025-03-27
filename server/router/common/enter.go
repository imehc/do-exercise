package common

import api "github.com/imehc/do-exercise/server/api/v1"

type RouterGroup struct {
	AuthRouter
}

var (
	authApi = api.ApiGroupApp.CommonApiGroup.AuthApi
)
