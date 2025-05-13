package system

import "github.com/gin-gonic/gin"

type SysInfoRouter struct{}

func (s *SysApiRouter) InitSysInfoRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("info")
	{
		router.GET("", sysInfoApi.Get) // 获取系统信息
	}
	return router
}
