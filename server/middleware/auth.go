package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/pkg/utils"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := global.OAUTH_SERVER.ValidationBearerToken(c.Request)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		var claims system.Claims
		username := token.GetUserID()
		if username == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		log := global.LOG
		db := global.DB
		cache := global.Cache
		var user system.User
		claimsCache, exist := cache.Get(global.CLAIMS)
		if exist {
			_ = json.Unmarshal([]byte(claimsCache.(string)), &claims)
			log.Info("Use cache")
		} else {
			err = db.Model(&system.User{}).Where("username = ?", username).First(&user).Error
			if err != nil {
				response.Unauthorized(c)
				c.Abort()
				return
			}
			claims.ID = user.ID
			claims.Username = user.Username
			claims.DeptId = user.DeptId
			claims.PostId = user.PostId
			claims.RoleId = user.RoleId

			expire, err := utils.ParseDurationString(global.CONFIG.Auth.AccessExpire)
			if err == nil {
				jsonData, err := json.Marshal(claims)
				if err == nil {
					cache.Set(global.CLAIMS, jsonData, expire)
				}
			}
			log.Info("Use database")
		}

		// TODO: 首先从redis中查出对应信息，如果redis中没有，则从数据库查出对应信息，存入上下文
		// TODO: 从redis中查出对应信息，比如租户ID,存入上下文, 住户解析验证租户是否存在，租户信息是否缺失
		c.Set(global.CLAIMS, claims)
		// c.Set("user_id", claims["user_id"])
		// c.Set("tenant_id", claims["tenant_id"])
		// c.Set("roles", claims["roles"])

		// 从上下文中取出用户角色，租户ID
		// userRole := username
		// tenantID := 1
		// // 构造 Casbin 请求：sub=角色, obj=路径, act=方法, tenant=租户
		// ok, err := global.Enforcer.Enforce(userRole, c.Request.URL.Path, c.Request.Method, fmt.Sprintf("tenant_%d", tenantID))
		// if err != nil || !ok {
		// 	response.Forbidden(c)
		// 	return
		// }

		c.Next()
	}
}
