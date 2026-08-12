package common

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/middleware"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("auth")
	{
		// 昂贵的无认证端点按 IP 单独限流：RSA 签发、验证码生成、邮件发送、登录/改密尝试。
		router.GET("code_with_email", middleware.IpRateLimitMiddleware(60*time.Second, 3), authApi.SendLoginWithEmailCode) // 获取使用邮箱登录验证码
		router.POST("login_with_email", middleware.IpRateLimitMiddleware(10*time.Second, 5), authApi.LoginWithEmail)       // 邮箱登录
		router.POST("login", middleware.IpRateLimitMiddleware(10*time.Second, 5), authApi.Login)                           // 登录
		router.POST("refresh_token", middleware.IpRateLimitMiddleware(10*time.Second, 5), authApi.RefreshToken)            // 刷新token
		router.GET("captcha", middleware.IpRateLimitMiddleware(10*time.Second, 5), authApi.GetCaptcha)                     // 获取验证码
		router.GET("public_key", middleware.IpRateLimitMiddleware(10*time.Second, 10), authApi.PublicKey)                  // 获取公钥
		router.PATCH("forget_password", middleware.IpRateLimitMiddleware(10*time.Second, 5), authApi.ResetPassword)        // 忘记密码
		router.GET("forget_password_code", middleware.IpRateLimitMiddleware(60*time.Second, 3), authApi.SendResetPasswordCode) // 发送忘记密码邮箱验证码
		router.Use(middleware.AuthMiddleware(), middleware.ContextMiddleware()).GET("logout", authApi.Logout) // 退出登录
	}
	return router
}
