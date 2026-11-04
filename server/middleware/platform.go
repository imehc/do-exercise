package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
)

// PlatformMiddleware protects platform operations independently of editable RBAC policies.
func PlatformMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("tenantId") != global.PlatformTenantID || !c.GetBool("isSuperAdmin") {
			response.Forbidden(c, "")
			c.Abort()
			return
		}
		c.Next()
	}
}
