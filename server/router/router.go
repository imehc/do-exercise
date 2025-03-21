package router

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/middleware"
)

func Run() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(middleware.ValidaterMiddleware())

	system := RouterGroupApp.System

	protected := r.Group("/system")
	public := r.Group("/")
	{
		// 健康监测
		public.GET("/health", func(c *gin.Context) {

		})
	}
	{
		system.InitSysUserRouter(protected)
	}

	return r
}
