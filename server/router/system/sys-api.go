package system

import "github.com/gin-gonic/gin"

type SysApiRouter struct{}

func (s *SysApiRouter) InitSysApiRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("apis")
	{
		router.PUT(":id", sysApiApi.Update)              // 更新api
		router.GET(":id", sysApiApi.Get)                 // 获取单个api
		router.GET("", sysApiApi.GetList)                // 获取api列表
		router.GET("all", sysApiApi.GetAll)              // 获取所有api
		router.GET("group_type", sysApiApi.GetGroupType) // 获取分组类型
	}
	return router
}
