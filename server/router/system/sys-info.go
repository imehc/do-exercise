package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/middleware"
)

type SysInfoRouter struct{}

func (s *SysInfoRouter) InitSysInfoRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("info", middleware.PlatformMiddleware())
	{
		router.GET("", sysInfoApi.Get) // 获取系统信息
	}
	return router
}
