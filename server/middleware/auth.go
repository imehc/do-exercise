package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/util"
	"net/http"
)

// mustChangePasswordWhitelist 强制改密期间仍允许访问的接口。
// 改密流程必须走通这几个端点，其余接口一律返回 403。
var mustChangePasswordWhitelist = map[string]bool{
	"PATCH /user/modify_password": true, // 修改密码
	"GET /auth/logout":            true, // 退出登录
	"GET /user/profile":           true, // 基本信息（前端布局依赖）
	"GET /user/menu":              true, // 菜单（前端布局依赖）
	"GET /sse":                    true, // SSE 连接（前端布局依赖）
	"GET /auth/my_tenants":        true, // 可用租户（顶栏当前租户展示依赖，未改密期间仅展示不切换）
}

// AuthMiddleware 登录验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		tokenData, err := util.ValidateAccessToken(c.Request.Context(), accessToken)
		if err != nil {
			if errors.Is(err, util.ErrAuthDependencyUnavailable) {
				c.Header("Retry-After", "1")
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 强制改密。默认管理员口令公开，首次登录后必须改密才能使用其他功能。
		if tokenData.MustChangePassword {
			key := c.Request.Method + " " + c.FullPath()
			if !mustChangePasswordWhitelist[key] {
				response.Forbidden(c, global.I18.Translate("passwordMustChange", c.GetString("lang")))
				c.Abort()
				return
			}
		}

		// 将用户ID存储在上下文中
		c.Set("userId", tokenData.UserId)
		c.Set("username", tokenData.Username)
		c.Set("tenantId", tokenData.TenantId)
		c.Set("authorizedTenantIds", tokenData.AuthorizedTenantIds)
		c.Set("accessToken", accessToken)
		c.Set("isSuperAdmin", tokenData.IsSuperAdmin)
		c.Next()
	}
}
