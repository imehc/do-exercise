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
	cleanupOnce    sync.Once                       // 清理协程只启动一次
)

// startCleanup 定期清理过期IP记录，防止requestInfoMap无限增长导致内存泄漏
func startCleanup() {
	cleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(timeWindow)
			defer ticker.Stop()
			for range ticker.C {
				threshold := time.Now().Add(-2 * timeWindow)
				mutex.Lock()
				for ip, info := range requestInfoMap {
					if info.LastAccessTime.Before(threshold) {
						delete(requestInfoMap, ip)
					}
				}
				mutex.Unlock()
			}
		}()
	})
}

// IpLimitMiddleware IP限流器
func IpLimitMiddleware(c *gin.Context) {
	startCleanup()

	ip := c.ClientIP()
	mutex.Lock()

	info, exists := requestInfoMap[ip]
	// 如果IP不存在，初始化并添加到map中
	if !exists {
		requestInfoMap[ip] = &RequestInfo{LastAccessTime: time.Now(), RequestNum: 1}
		mutex.Unlock()
		c.Next()
		return
	}
	// 如果超过时间窗口，重置请求计数
	if time.Since(info.LastAccessTime) > timeWindow {
		info.RequestNum = 1
		info.LastAccessTime = time.Now()
		mutex.Unlock()
		c.Next()
		return
	}
	// 在时间窗口内，增加请求计数
	info.RequestNum++
	// 如果请求计数超过限制，禁止访问
	if info.RequestNum > maxRequests {
		mutex.Unlock()
		response.StatusTooManyRequests(c)
		c.Abort()
		return
	}
	// 更新最后访问时间
	info.LastAccessTime = time.Now()
	mutex.Unlock()
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
