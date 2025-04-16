package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
)

// CasbinHandler 基于Casbin的权限验证中间件
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的URI
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method

		// 从上下文中获取用户信息和域信息
		sub := c.GetString("userId")
		dom := c.GetString("domain")
		if sub == "" || dom == "" {
			response.Forbidden(c)
			c.Abort()
			return
		}

		// 检查权限，使用带有域的验证
		if ok, err := global.Enforcer.Enforce(sub, dom, obj, act); err != nil || !ok {
			response.Forbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
