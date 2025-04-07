package common

import (
	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("")
	{
		router.POST("login", authApi.Login)         // 登录
		router.GET("public_key", authApi.PublicKey) // 获取公钥
	}
	return router
}
