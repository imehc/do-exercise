package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
)

// ContextMiddleware 设置登录后全局上下文
func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.MustGet("userId").(string)
		ctx := context.WithValue(c.Request.Context(), global.ContextUserIDKey, userId)
		if tenantId := c.GetString("tenantId"); tenantId != "" {
			ctx = context.WithValue(ctx, global.ContextTenantIDKey, tenantId)
		}
		// 透传超管标识：AuthMiddleware 已校验过会话，服务层据此放开跨租户可见范围，
		// 不必在每个列表查询里再回查一次 sys_user.is_super_admin。
		if c.GetBool("isSuperAdmin") {
			ctx = context.WithValue(ctx, global.ContextIsSuperAdminKey, true)
		}
		// 请求范围内派生 DB，避免并发写全局 global.DB 造成数据竞争
		c.Set(global.ContextDBKey, global.DB.WithContext(ctx))
		c.Next()
	}
}

// RequestContextMiddleware 设置请求上下文
func RequestContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(global.ContextDBKey, global.DB.WithContext(c.Request.Context()))
		c.Next()
	}
}
