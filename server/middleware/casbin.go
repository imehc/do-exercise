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

		if ok, err := global.Enforcer.Enforce(sub, obj, act); err != nil || !ok {
			response.Forbidden(c, "")
			c.Abort()
			return
		}

		c.Next()
	}
}
