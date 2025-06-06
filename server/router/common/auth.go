package common

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/middleware"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("auth")
	{

		router.GET("code_with_email", authApi.SendLoginWithEmailCode)                                         // 获取使用邮箱登录验证码
		router.POST("login_with_email", authApi.LoginWithEmail)                                               // 邮箱登录
		router.POST("login", authApi.Login)                                                                   // 登录
		router.GET("refresh_token", authApi.RefreshToken)                                                     // 刷新token
		router.GET("captcha", authApi.GetCaptcha)                                                             // 获取验证码
		router.GET("public_key", authApi.PublicKey)                                                           // 获取公钥
		router.PATCH("forget_password", authApi.ResetPassword)                                                // 忘记密码
		router.GET("forget_password_code", authApi.SendResetPasswordCode)                                     // 发送忘记密码邮箱验证码
		router.Use(middleware.AuthMiddleware(), middleware.ContextMiddleware()).GET("logout", authApi.Logout) // 退出登录
	}
	return router
}
