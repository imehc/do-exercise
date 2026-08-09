//go:build windows
// +build windows

package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:    address,
		Handler: router,
		// 与非 Windows 分支保持一致：限制慢速连接占用 worker 的时长
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		// WriteTimeout 保持较长以兼容 SSE 长连接
		WriteTimeout: 10 * time.Minute,
		// 空闲的 keep-alive 连接 60s 未复用即回收，避免连接堆积
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
