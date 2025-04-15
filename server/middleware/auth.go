package middleware

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
)

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
		// TODO: 使用结构体解析
		var tokenData map[string]string
		if err := json.Unmarshal([]byte(tokenInfo), &tokenData); err != nil {
			response.ServerError(c)
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中
		c.Set("userId", tokenData["userId"])
		c.Next()
	}
}
