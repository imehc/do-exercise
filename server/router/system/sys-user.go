package system

import (
	"github.com/gin-gonic/gin"
)

type SysUserRouter struct{}

func (s *SysUserRouter) InitSysUserRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("")
	{
		r1 := r.Group("users")
		r1.POST("", sysUserApi.Create) // 创建用户
	}
	{
		router.POST("login", sysUserApi.Login) // 登录
	}
	return router
}
