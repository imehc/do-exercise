//go:build !windows
// +build !windows

package core

import (
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	// 读取请求头 5s、读完整个请求 30s：原先仅设 ReadHeaderTimeout=10min，
	// 等同于没有限制——几百个慢速连接各占一个 worker 十分钟即可打满服务（slowloris）。
	s.ReadHeaderTimeout = 5 * time.Second
	s.ReadTimeout = 30 * time.Second
	// WriteTimeout 保持较长：SSE 等长连接依赖它。
	// 若后续按审计建议摘掉或改造 SSE，可下调至 30s~60s。
	s.WriteTimeout = 10 * time.Minute
	s.MaxHeaderBytes = 1 << 20
	// endless 自身处理 SIGINT/SIGTERM/SIGHUP 并等待在途请求结束，故此处无需再接信号
	return s
}
