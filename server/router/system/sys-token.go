package system

import (
	"github.com/gin-gonic/gin"
)

type SysTokenRouter struct{}

func (s *SysTokenRouter) InitSysTokenRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("token")
	{
		router.GET("all", sysTokenApi.FindAll)     // 获取token列表
		router.DELETE("", sysTokenApi.Delete)      // 删除token
		router.PATCH("", sysTokenApi.ModityStatus) // 修改token状态
	}
	return router
}
