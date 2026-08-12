package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
)

type RequestInfo struct {
	LastAccessTime time.Time // 上次访问时间
	RequestNum     int       // 请求计数
}

// 全局兜底限流参数：所有请求默认 10 次/秒/IP。
const (
	maxRequests = 10
	timeWindow  = 1 * time.Second
)

// ipLimiter 按 IP 维度做滑动窗口限流。
// 相比全局 token bucket，按 IP 隔离不会误伤正常的多用户并发流量。
type ipLimiter struct {
	mutex       sync.Mutex
	reqs        map[string]*RequestInfo
	timeWindow  time.Duration
	maxRequests int
}

func newIPLimiter(timeWindow time.Duration, maxRequests int) *ipLimiter {
	l := &ipLimiter{
		reqs:        make(map[string]*RequestInfo),
		timeWindow:  timeWindow,
		maxRequests: maxRequests,
	}
	go l.startCleanup()
	return l
}

// startCleanup 定期清理过期的 IP 记录，防止 map 无限增长导致内存泄漏。
func (l *ipLimiter) startCleanup() {
	ticker := time.NewTicker(l.timeWindow)
	defer ticker.Stop()
	for range ticker.C {
		threshold := time.Now().Add(-2 * l.timeWindow)
		l.mutex.Lock()
		for ip, info := range l.reqs {
			if info.LastAccessTime.Before(threshold) {
				delete(l.reqs, ip)
			}
		}
		l.mutex.Unlock()
	}
}

// Handle 限流处理。锁只在状态更新期间持有，c.Next() 在解锁后执行，避免长请求持锁阻塞其它 IP。
func (l *ipLimiter) Handle(c *gin.Context) {
	ip := c.ClientIP()
	l.mutex.Lock()

	info, exists := l.reqs[ip]
	// 如果IP不存在，初始化并添加到map中
	if !exists {
		l.reqs[ip] = &RequestInfo{LastAccessTime: time.Now(), RequestNum: 1}
		l.mutex.Unlock()
		c.Next()
		return
	}
	// 如果超过时间窗口，重置请求计数
	if time.Since(info.LastAccessTime) > l.timeWindow {
		info.RequestNum = 1
		info.LastAccessTime = time.Now()
		l.mutex.Unlock()
		c.Next()
		return
	}
	// 在时间窗口内，增加请求计数
	info.RequestNum++
	// 如果请求计数超过限制，禁止访问
	if info.RequestNum > l.maxRequests {
		l.mutex.Unlock()
		response.StatusTooManyRequests(c)
		c.Abort()
		return
	}
	// 更新最后访问时间
	info.LastAccessTime = time.Now()
	l.mutex.Unlock()
	c.Next()
}

// IpLimitMiddleware IP限流器（全局兜底：所有请求 10 次/秒/IP）
var IpLimitMiddleware = newIPLimiter(timeWindow, maxRequests).Handle

// IpRateLimitMiddleware 生成按 IP 限流的中间件，用于对昂贵端点单独收紧配额。
// 例：验证码生成、RSA 密钥签发、邮件发送。
func IpRateLimitMiddleware(timeWindow time.Duration, maxRequests int) gin.HandlerFunc {
	return newIPLimiter(timeWindow, maxRequests).Handle
}