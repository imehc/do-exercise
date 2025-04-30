package system

import (
	"github.com/gin-gonic/gin"
)

type SysOperationLogRouter struct{}

func (s *SysOperationLogRouter) InitSysOperationLogRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("logs")
	{
		router.GET("", sysOperationLogApi.GetList) // 获取操作记录列表
	}
	return router
}
