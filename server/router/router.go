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
	r.Use(middleware.IpLimitMiddleware)

	system := RouterGroupApp.System
	common := RouterGroupApp.Common

	protected := r.Group("/system")
	protected.Use(
		middleware.AuthMiddleware(),
		middleware.ContextMiddleware(),
		middleware.CasbinMiddleware(),
		middleware.OperationLogMiddleware(),
	)

	noAuth := r.Group("/")
	noAuth.Use(
		middleware.OperationLogMiddleware(),
	)

	auth := r.Group("/")
	auth.Use(
		middleware.AuthMiddleware(),
		middleware.ContextMiddleware(),
		middleware.OperationLogMiddleware(),
	)

	public := r.Group("/")
	{
		// 健康监测
		public.GET("/health", func(c *gin.Context) {

		})
	}
	{
		system.InitSysApiRouter(protected)
		system.InitSysUserRouter(protected)
		system.InitSysMenuRouter(protected)
		system.InitSysRoleRouter(protected)
		system.InitSysOperationLogRouter(protected)
		common.InitAuthRouter(noAuth)
		common.InitUserRouter(auth)
	}

	return r
}
