package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/spf13/cast"
)

// CasbinMiddleware 基于Casbin的权限验证中间件
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Request.URL.Path
		act := c.Request.Method

		userIdString, exists := c.Get("userId")
		sub := cast.ToString(userIdString)
		if !exists || sub == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 多租户模式下以租户ID为 Casbin 域；单租户为默认租户。
		dom := c.GetString("tenantId")
		if dom == "" {
			dom = global.Config.Tenant.DefaultTenantId
		}

		// 平台超级管理员由 is_super_admin 标识直接放行，不再依赖播种的平台角色/Casbin 策略。
		// 仅平台域（dom=platform）的账号能被标记为平台超管，业务租户永不生效。
		if dom == global.PlatformTenantID && c.GetBool("isSuperAdmin") {
			c.Next()
			return
		}

		if ok, err := global.Enforcer.Enforce(sub, dom, obj, act); err != nil || !ok {
			response.Forbidden(c, "")
			c.Abort()
			return
		}

		c.Next()
	}
}
