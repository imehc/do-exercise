package router

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/middleware"
)

func Run() *gin.Engine {
	internal.InitCaptcha()
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(middleware.ValidaterMiddleware())

	system := RouterGroupApp.System
	common := RouterGroupApp.Common

	protected := r.Group("/system")
	protected.Use(middleware.AuthMiddleware())
	basic := r.Group("/")
	public := r.Group("/")
	{
		// 健康监测
		public.GET("/health", func(c *gin.Context) {

		})
	}
	{
		system.InitSysUserRouter(protected)
		common.InitAuthRouter(basic)
	}

	return r
}
