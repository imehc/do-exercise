package system

import (
	"github.com/gin-gonic/gin"
)

type SysUserRouter struct{}

func (s *SysUserRouter) InitSysUserRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("users")
	{
		router.POST("", sysUserApi.Create) // 创建用户
	}
	return router
}
