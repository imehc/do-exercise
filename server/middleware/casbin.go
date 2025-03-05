package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
)

func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := c.Get(global.CLAIMS)
		claims := data.(system.Claims)

		// 如果是admin角色，直接放行
		if claims.IsAdmin {
			c.Next()
			return
		}

		// 获取请求的URI
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := claims.RoleId

		// 检查权限
		if success, _ := global.Enforcer.Enforce(sub, obj, act); !success {
			response.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
