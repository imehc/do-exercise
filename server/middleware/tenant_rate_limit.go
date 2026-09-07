package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
)

// 租户聚合限流的回退默认值。配置缺失或写成 0 时兜底，避免一次配置缺省让整层失效。
const (
	tenantRLWindowSeconds  = 1
	tenantRLBusinessPerSec = 200
	tenantRLPlatformPerSec = 1000
)

// TenantRateLimitMiddleware 租户级聚合请求限流。
// 挂在鉴权/租户上下文解析之后：按当前租户给一段固定窗口内的聚合请求预算，
// 单个租户突发流量打满自己的预算即被拒绝，不影响其它租户与平台域。
// 平台域（超管、无租户上下文）用独立且更高的预算档。
func TenantRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := c.GetString("tenantId")
		if tenantRateLimitBlocked(c.Request.Context(), tid) {
			response.StatusTooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// tenantRateLimitBlocked 判断当前租户在本次固定窗口内是否已超预算。
// tenantId 为空表示平台域请求（超管等），使用 platform 档。
// Redis 不可用或计数出错时按“未超限”放行，避免限流本身拖垮请求。
func tenantRateLimitBlocked(ctx context.Context, tenantId string) bool {
	rl := global.Config.RateLimit
	if !rl.Enable {
		return false
	}
	window := rl.WindowSeconds
	if window <= 0 {
		window = tenantRLWindowSeconds
	}

	effectiveID := tenantId
	perSec := rl.BusinessPerSec
	if effectiveID == "" {
		effectiveID = global.PlatformTenantID
		perSec = rl.PlatformPerSec
	}
	if perSec <= 0 {
		if effectiveID == global.PlatformTenantID {
			perSec = tenantRLPlatformPerSec
		} else {
			perSec = tenantRLBusinessPerSec
		}
	}
	budget := perSec * window

	slot := time.Now().Unix() / int64(window)
	key := fmt.Sprintf("rl:tenant:%s:%d", effectiveID, slot)

	val, err := global.Redis.Incr(ctx, key).Result()
	if err != nil {
		return false
	}
	if val == 1 {
		// 只在新窗口第一笔时设过期，避免每笔都打 EXPIRE
		global.Redis.Expire(ctx, key, time.Duration(window)*time.Second).Result()
	}
	return val > int64(budget)
}
