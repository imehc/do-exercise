package middleware

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/system"
	sysSevice "github.com/imehc/do-exercise/server/service/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/mssola/user_agent"
)

// 操作日志中请求体/响应体的最大存储长度，防止大请求撑爆日志表
const maxStoredBodyLength = 2 * 1024

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

		// 包装 ResponseWriter 以获取响应内容
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

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
			Params:         c.Request.URL.RawQuery,
			Body:           truncateString(string(body), maxStoredBodyLength),
			Result:         truncateString(blw.body.String(), maxStoredBodyLength),
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
			logService.WithContext(ctx).Create(operationLog)
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
