package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
)

func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		enforce := global.Enforcer
		// 获取请求的URI
		obj := c.Request.URL.Path
		// 获取请求方法
		act := c.Request.Method

		data, _ := c.Get(global.CLAIMS)
		claims := data.(system.Claims)

		sub := claims.Username
		dom := strconv.Itoa(claims.DeptId)

		// 检查权限
		ok, err := enforce.Enforce(
			sub, // sub
			dom, // dom
			obj, // ibj
			act, // act
		)
		if err != nil || !ok {
			response.Forbidden(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
