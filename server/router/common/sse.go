package common

import (
	"io"

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
	c.Header("Access-Control-Allow-Origin", "*")

	// 创建新的客户端连接
	client := s.Manager.AddClient(userId)
	defer s.Manager.RemoveClient(userId)

	c.Stream(func(w io.Writer) bool {
		// 从客户端通道接收消息
		if msg, ok := <-client.Channel; ok {
			c.SSEvent("message", string(msg))
			return true
		}
		return false
	})
}
