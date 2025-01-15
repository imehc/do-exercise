package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/middleware"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	r := Router.Group("")
	r2 := r.Group("").Use(middleware.ResponseError())
	{
		r.POST("login", authApi.Login)
		r.GET("captcha", captchaApi.GetCaptcha)
	}
	{
		r2.GET("refresh_token", authApi.RefreshToken)
		r2.GET("logout", authApi.Logout)
	}
	return r
}
