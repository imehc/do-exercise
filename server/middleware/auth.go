package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := global.OAUTH_SERVER.ValidationBearerToken(c.Request)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		claims := system.Claims{}
		username := token.GetUserID()
		if username != "" {
			claims.Username = username
		}
		// TODO: 首先从redis中查出对应信息，如果redis中没有，则从数据库查出对应信息，存入上下文
		c.Set("claims", claims)
		c.Next()
	}
}
