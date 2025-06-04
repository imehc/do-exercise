package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/juju/ratelimit"
)

type RequestInfo struct {
	LastAccessTime time.Time // 上次访问时间
	RequestNum     int       // 请求计数
}

var (
	requestInfoMap = make(map[string]*RequestInfo) // IP到请求信息的映射
	mutex          = &sync.Mutex{}                 // 用于保护requestInfoMap的互斥锁
	maxRequests    = 10                            // 允许的最大请求数
	timeWindow     = 1 * time.Second               // 时间窗口
)

// IpLimitMiddleware IP限流器
func IpLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()
	mutex.Lock()
	defer mutex.Unlock()
	// 检查IP是否在map中
	info, exists := requestInfoMap[ip]
	// 如果IP不存在，初始化并添加到map中
	if !exists {
		requestInfoMap[ip] = &RequestInfo{LastAccessTime: time.Now(), RequestNum: 1}
		return
	}
	// 如果IP存在，检查时间窗口
	if time.Since(info.LastAccessTime) > timeWindow {
		// 如果超过时间窗口，重置请求计数
		info.RequestNum = 1
		info.LastAccessTime = time.Now()
		return
	}
	info.RequestNum++ // 如果在时间窗口内，增加请求计数
	// 如果请求计数超过限制，禁止访问
	if info.RequestNum > maxRequests {
		response.StatusTooManyRequests(c)
		c.Abort()
		return
	}
	// 更新最后访问时间
	info.LastAccessTime = time.Now()
	c.Next()
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(time time.Duration, originNum, pushNum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(time, originNum, pushNum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.StatusTooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
