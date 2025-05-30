package router

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/middleware"
	"github.com/imehc/do-exercise/server/model/common/response"
	sse "github.com/imehc/do-exercise/server/router/common"
	"github.com/imehc/do-exercise/server/util"
)

func Run() *gin.Engine {
	if util.IsRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	internal.InitCaptcha()
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestContextMiddleware())
	r.Use(middleware.OperationLogMiddleware())
	r.Use(middleware.ValidaterMiddleware())
	r.Use(middleware.IpLimitMiddleware)

	co := sse.NewSSERouter()
	shared.SSEManager = co.Manager

	system := RouterGroupApp.System
	common := RouterGroupApp.Common

	protected := r.Group("/system")
	protected.Use(
		middleware.AuthMiddleware(),
		middleware.ContextMiddleware(),
		middleware.CasbinMiddleware(),
	)

	noAuth := r.Group("/")

	auth := r.Group("/")
	auth.Use(
		middleware.AuthMiddleware(),
		middleware.ContextMiddleware(),
	)

	public := r.Group("/")
	{
		// 健康监测
		public.GET("/health", func(c *gin.Context) {
			response.Success(c, "ok")
		})
	}
	{
		system.InitSysApiRouter(protected)
		system.InitSysUserRouter(protected)
		system.InitSysMenuRouter(protected)
		system.InitSysRoleRouter(protected)
		system.InitSysOperationLogRouter(protected)
		system.InitSysTokenRouter(protected)
		system.InitSysInfoRouter(protected)

		common.InitAuthRouter(noAuth)

		co.InitSSERouter(auth)
		common.InitUserRouter(auth)
		common.InitOssRouter(auth)
	}

	return r
}
