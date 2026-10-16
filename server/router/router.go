package router

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/middleware"
	"github.com/imehc/do-exercise/server/model/common/response"
	sse "github.com/imehc/do-exercise/server/router/common"
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
)

// Run 构建路由并初始化验证码，供服务启动与集成测试使用。
func Run() *gin.Engine {
	if util.IsRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	internal.InitCaptcha()
	return New()
}

// New 纯构建路由，不依赖 Redis/验证码等运行时依赖，供 openapi 漂移检查等离线工具复用。
// 注意：InitCaptcha 在 Run 中完成，New 调用方如需验证码请自行初始化。
func New() *gin.Engine {
	r := gin.Default()

	// 限定受信任的反向代理。
	// gin 默认信任所有代理，此时任何客户端都能通过伪造 X-Forwarded-For 决定
	// c.ClientIP() 的返回值，从而绕过 IP 限流与登录锁定，并污染审计日志。
	if err := r.SetTrustedProxies(global.Config.System.TrustedProxies); err != nil {
		global.Log.Error("配置受信任代理失败，将不信任任何代理",
			zap.Strings("trustedProxies", global.Config.System.TrustedProxies),
			zap.Error(err))
		_ = r.SetTrustedProxies(nil)
	}

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
		middleware.TenantRateLimitMiddleware(),
		middleware.CasbinMiddleware(),
	)

	noAuth := r.Group("/")

	// 认证但不鉴权的系统组：数据是全租户共用的字典，任何登录会话都可读。
	// 见 InitSysMenuTreeRouter 的说明（P2-2）。
	authedSystem := r.Group("/system")
	authedSystem.Use(
		middleware.AuthMiddleware(),
		middleware.ContextMiddleware(),
		middleware.TenantRateLimitMiddleware(),
	)

	auth := r.Group("/")
	auth.Use(
		middleware.AuthMiddleware(),
		middleware.ContextMiddleware(),
		middleware.TenantRateLimitMiddleware(),
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
		system.InitSysMenuTreeRouter(authedSystem)
		system.InitSysRoleRouter(protected)
		system.InitSysOperationLogRouter(protected)
		system.InitSysTokenRouter(protected)
		system.InitSysInfoRouter(protected)
		system.InitSysJobRouter(protected)
		system.InitSysTenantRouter(protected)

		common.InitAuthRouter(noAuth)

		co.InitSSERouter(auth)
		common.InitUserRouter(auth)
		common.InitOssRouter(auth)
	}

	return r
}
