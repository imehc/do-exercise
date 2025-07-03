package system

import "github.com/gin-gonic/gin"

type SysJobRouter struct{}

func (s *SysJobRouter) InitSysJobRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("jobs")
	{
		router.POST("", sysJobApi.Create)             // 创建定时任务
		router.DELETE(":id", sysJobApi.Delete)        // 删除定时任务
		router.PUT(":id", sysJobApi.Update)           // 更新定时任务
		router.GET(":id", sysJobApi.Get)              // 获取单个定时任务
		router.GET("", sysJobApi.GetList)             // 获取定时任务列表
		router.POST(":id/start", sysJobApi.Start)     // 启动任务
		router.POST(":id/stop", sysJobApi.Stop)       // 停止任务
		router.POST(":id/execute", sysJobApi.Execute) // 立即执行一次
	}
	return router
}
