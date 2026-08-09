package middleware

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common/response"
)

// mustChangePasswordWhitelist 强制改密期间仍允许访问的接口。
// 改密流程必须走通这几个端点，其余接口一律返回 403。
var mustChangePasswordWhitelist = map[string]bool{
	"PATCH /user/modify_password": true, // 修改密码
	"GET /auth/logout":            true, // 退出登录
	"GET /user/profile":           true, // 基本信息（前端布局依赖）
	"GET /user/menu":              true, // 菜单（前端布局依赖）
	"GET /sse":                    true, // SSE 连接（前端布局依赖）
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

		// 验证token
		ctx := context.Background()
		tokenInfo, err := global.Redis.Get(ctx, fmt.Sprintf("accessToken:%s", accessToken)).Result()
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 解析token信息
		var tokenData model.TokenInfo
		if err := json.Unmarshal([]byte(tokenInfo), &tokenData); err != nil {
			response.ServerError(c)
			c.Abort()
			return
		}

		if tokenData.Disabled {
			response.Forbidden(c, global.I18.Translate("accessTokenIsDisabled", c.GetString("lang")))
			c.Abort()
			return
		}

		// 强制改密。默认管理员口令公开，首次登录后必须改密才能使用其他功能。
		if tokenData.MustChangePassword {
			key := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
			if !mustChangePasswordWhitelist[key] {
				response.Forbidden(c, global.I18.Translate("passwordMustChange", c.GetString("lang")))
				c.Abort()
				return
			}
		}

		// 将用户ID存储在上下文中
		c.Set("userId", tokenData.UserId)
		c.Set("roles", tokenData.RoleIds)
		c.Set("username", tokenData.Username)
		c.Set("accessToken", accessToken)
		c.Next()
	}
}
