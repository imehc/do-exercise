package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/middleware"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("")

	r2 := r.Group("").Use(middleware.JWTAuth())
	{
		r.GET("captcha", captchaApi.GetCaptcha)
		r.POST("login", authApi.Login)
		r.GET("refresh_token", authApi.RefreshToken)
	}
	{
		r2.GET("logout", authApi.Logout)
	}
	return r
}
