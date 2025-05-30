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
		global.DB = global.DB.WithContext(ctx)

		c.Next()
	}
}

// RequestContextMiddleware 设置请求上下文
func RequestContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.DB = global.DB.WithContext(c.Request.Context())
		c.Next()
	}
}
