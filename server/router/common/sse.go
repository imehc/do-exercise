package common

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/cast"
)

type SSERouter struct {
	Manager *util.Manager
}

// NewSSERouter 创建新的SSERouter实例
func NewSSERouter() *SSERouter {
	return &SSERouter{
		Manager: util.NewManager(),
	}
}

// InitSSERouter 初始化SSE路由
func (s *SSERouter) InitSSERouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("/sse")
	{
		// SSE连接端点
		router.GET("", s.handleSSE)
	}
	return router
}

// sseHeartbeatInterval 心跳间隔。
// nginx 未显式配置 proxy_read_timeout 时默认 60s 会切断空闲连接，
// 心跳必须显著短于该值。
const sseHeartbeatInterval = 20 * time.Second

// handleSSE 处理SSE连接请求
func (s *SSERouter) handleSSE(c *gin.Context) {
	userIdString, exists := c.Get("userId")
	userId := cast.ToString(userIdString)
	if !exists || userId == "" {
		response.Unauthorized(c)
		return
	}

	// 设置SSE相关的响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// 关闭 nginx 缓冲，否则消息会被攒着不下发
	c.Header("X-Accel-Buffering", "no")

	// 每条连接独立标识，支持同一用户多标签页/多设备
	connID, err := util.Uuid()
	if err != nil {
		response.ServerError(c)
		return
	}

	client := s.Manager.AddClient(userId, connID)
	defer s.Manager.RemoveClient(userId, connID)

	ticker := time.NewTicker(sseHeartbeatInterval)
	defer ticker.Stop()

	ctx := c.Request.Context()

	// 关键：step 内必须 select 上 ctx.Done()。
	// gin 的 Stream 只在两次 step 调用之间检查客户端断开，
	// 若 step 阻塞在 <-client.Channel，断开信号永远读不到，
	// handler goroutine 会永久 park 且 defer 不执行——每次断连泄漏一个 goroutine。
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client.Channel:
			if !ok {
				return false
			}
			c.SSEvent("message", string(msg))
			return true
		case <-ticker.C:
			// 注释帧：保活，同时让 Stream 有机会回到外层检查连接状态
			c.SSEvent("ping", "")
			return true
		case <-ctx.Done():
			return false
		}
	})
}
