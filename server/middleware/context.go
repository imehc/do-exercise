package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
)

// ContextHandler 设置全局上下文
func ContextHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.MustGet("userId").(int64)
		ctx := context.WithValue(c.Request.Context(), "userId", userId)
		global.DB = global.DB.WithContext(ctx)

		c.Next()
		return
	}
}
