package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/util"
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

		// Casbin 域恒等于租户 ID。会话没有租户是不该出现的状态（登录一定会落一个
		// 业务租户或 platform），这里 fail-closed，不再兜底到某个具体租户——
		// 兜底会让一个残缺会话拿到该租户的全部策略。
		dom := c.GetString("tenantId")
		if dom == "" {
			response.Forbidden(c, "")
			c.Abort()
			return
		}

		// 平台超级管理员由 is_super_admin 标识直接放行，不再依赖播种的平台角色/Casbin 策略。
		// 仅平台域（dom=platform）的账号能被标记为平台超管，业务租户永不生效。
		if dom == global.PlatformTenantID && c.GetBool("isSuperAdmin") {
			c.Next()
			return
		}

		if ok, err := util.Authorize(util.DB(c), sub, dom, obj, act); err != nil || !ok {
			response.Forbidden(c, "")
			c.Abort()
			return
		}

		c.Next()
	}
}
