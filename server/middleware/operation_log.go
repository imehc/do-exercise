package middleware

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	sysSevice "github.com/imehc/do-exercise/server/service/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/mssola/user_agent"
	"go.uber.org/zap"
)

// 操作日志中请求体/响应体的最大存储长度，防止大请求撑爆日志表
const maxStoredBodyLength = 2 * 1024

// 敏感字段占位符
const redactedPlaceholder = "[REDACTED]"

// isSensitivePath 判断该路径的响应体/查询串是否含凭据。
// /auth/* 的响应体是 access/refresh token，refresh_token 还会出现在查询串里，
// 一旦落库，任何拥有日志读取权限的人都能直接重放会话。
func isSensitivePath(path string) bool {
	return strings.HasPrefix(path, "/auth/")
}

// skipResponseCapture 判断是否跳过响应体缓冲。
// 除敏感路径外，SSE 是长连接，响应永不结束，缓冲区会随连接生命周期无限增长。
func skipResponseCapture(path string) bool {
	return isSensitivePath(path) || path == "/sse"
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 获取请求体
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// 包装 ResponseWriter 以获取响应内容。
		// 敏感路径与 SSE 不缓冲：前者会把凭据写进日志表，后者会无限增长。
		captureResponse := !skipResponseCapture(c.Request.URL.Path)
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		if captureResponse {
			c.Writer = blw
		}

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 获取客户端信息
		uaRaw := c.GetHeader("User-Agent")
		ua := user_agent.New(uaRaw)
		browserName, browserVersion := ua.Browser()
		os := ua.OS()

		clientIP := c.ClientIP()
		isInternal := isInternalIP(clientIP)

		// 创建操作日志
		sensitive := isSensitivePath(c.Request.URL.Path)

		// 查询串在 /auth/refresh_token 上携带 refresh_token，必须脱敏
		params := c.Request.URL.RawQuery
		if sensitive && params != "" {
			params = redactedPlaceholder
		}

		// 响应体在 /auth/login 等路径上就是 access/refresh token 本身。
		// 请求体保留：登录口令由前端 RSA 加密后传输，且保留用户名对失败登录的溯源有价值。
		result := redactedPlaceholder
		if captureResponse {
			result = truncateString(blw.body.String(), maxStoredBodyLength)
		}

		operationLog := system.SysOperationLog{
			Username:       c.GetString("username"),
			Ip:             clientIP,
			IsInternalIP:   isInternal,
			UserAgent:      uaRaw,
			Borwser:        browserName,
			BorwserVersion: browserVersion,
			IsMobile:       ua.Mobile(),
			IsBot:          ua.Bot(),
			Os:             os,
			Method:         c.Request.Method,
			Path:           c.Request.URL.Path,
			Params:         params,
			Body:           truncateString(string(body), maxStoredBodyLength),
			Result:         result,
			Success:        c.Writer.Status() >= 200 && c.Writer.Status() < 400,
			Code:           c.Writer.Status(),
			Message:        c.Errors.String(), // 获取错误信息
			StartTime:      startTime,
			EndTime:        endTime,
			Latency:        endTime.Sub(startTime).Milliseconds(),
		}

		if !isInternal {
			operationLog.Address = util.IPToRegion(clientIP)
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var logService sysSevice.SysOperationLogService
			if err := logService.WithContext(ctx).Create(operationLog); err != nil {
				global.Log.Error("写入操作日志失败",
					zap.String("path", operationLog.Path),
					zap.String("method", operationLog.Method),
					zap.Error(err))
			}
		}()
	}
}

func truncateString(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

// 判断是否为内网IP
func isInternalIP(ip string) bool {
	// 本地回环地址
	if ip == "127.0.0.1" || ip == "::1" || ip == "localhost" {
		return true
	}

	// 内网IP地址段
	privateIPRanges := []struct {
		start string
		end   string
	}{
		{"10.0.0.0", "10.255.255.255"},                        // 10.0.0.0/8
		{"172.16.0.0", "172.31.255.255"},                      // 172.16.0.0/12
		{"192.168.0.0", "192.168.255.255"},                    // 192.168.0.0/16
		{"169.254.0.0", "169.254.255.255"},                    // 169.254.0.0/16
		{"fc00::", "fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"}, // ULA
	}

	// 将IP地址转换为整数进行比较
	ipBytes := net.ParseIP(ip)
	if ipBytes == nil {
		return false
	}

	if ipBytes.To4() != nil {
		// IPv4
		ipInt := binary.BigEndian.Uint32(ipBytes.To4())
		for _, r := range privateIPRanges[:4] { // 只检查IPv4范围
			startInt := binary.BigEndian.Uint32(net.ParseIP(r.start).To4())
			endInt := binary.BigEndian.Uint32(net.ParseIP(r.end).To4())
			if ipInt >= startInt && ipInt <= endInt {
				return true
			}
		}
	} else {
		// IPv6
		if ipBytes.IsLoopback() || ipBytes.IsLinkLocalUnicast() || ipBytes.IsPrivate() {
			return true
		}
	}

	return false
}
