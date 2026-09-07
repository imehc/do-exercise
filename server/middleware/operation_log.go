package middleware

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/mssola/user_agent"
	"go.uber.org/zap"
)

// 操作日志中请求体/响应体的最大存储长度，防止大请求撑爆日志表
const maxStoredBodyLength = 2 * 1024

// 敏感字段占位符
const redactedPlaceholder = "[REDACTED]"

// 批量落库参数：队列容量、单批条数、最大攒批时长。
// 有界队列保证高并发下不会无界地开 goroutine 写库把连接池占满。
const (
	operationLogQueueSize   = 1024
	operationLogBatchSize   = 100
	operationLogFlushPeriod = 1 * time.Second
)

var (
	operationLogQueue      = make(chan system.SysOperationLog, operationLogQueueSize)
	operationLogWorkerOnce sync.Once
)

// StartOperationLogWorker 启动操作日志批量落库 worker，进程启动时调用一次。
func StartOperationLogWorker() {
	operationLogWorkerOnce.Do(func() {
		go operationLogWorker()
	})
}

// operationLogWorker 从有界队列取日志，攒满一批或超时后批量 INSERT。
func operationLogWorker() {
	batch := make([]system.SysOperationLog, 0, operationLogBatchSize)
	ticker := time.NewTicker(operationLogFlushPeriod)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		logs := batch
		batch = make([]system.SysOperationLog, 0, operationLogBatchSize)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := global.DB.WithContext(ctx).Create(&logs).Error; err != nil {
			global.Log.Error("批量写入操作日志失败",
				zap.Int("count", len(logs)),
				zap.Error(err))
		}
	}

	for {
		select {
		case log := <-operationLogQueue:
			batch = append(batch, log)
			if len(batch) >= operationLogBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

// enqueueOperationLog 非阻塞入队。队列满时丢弃该条日志并告警，
// 避免日志写入反过来拖垮业务请求。
func enqueueOperationLog(log system.SysOperationLog) {
	select {
	case operationLogQueue <- log:
	default:
		global.Log.Warn("操作日志队列已满，丢弃日志",
			zap.String("path", log.Path),
			zap.String("method", log.Method))
	}
}

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

		// 获取请求体。
		// 只在 body 体积可预知且不超过上限时才缓冲（供日志截取）：
		// 超大 POST、上传（multipart）、未知长度的分块 body 一律不碰，
		// 原样留给 handler，避免把整个请求体读进内存。
		var body []byte
		contentLength := c.Request.ContentLength
		if contentLength > 0 && contentLength <= maxStoredBodyLength &&
			!strings.HasPrefix(c.ContentType(), "multipart/form-data") {
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
		// 未缓冲 body 的请求（大/上传/未知长度）用占位符表示。
		result := redactedPlaceholder
		if captureResponse {
			result = truncateString(blw.body.String(), maxStoredBodyLength)
		}

		requestBody := redactedPlaceholder
		if body != nil {
			requestBody = truncateString(string(body), maxStoredBodyLength)
		}

		operationLog := system.SysOperationLog{
			Username:       c.GetString("username"),
			TenantId:       c.GetString("tenantId"),
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
			Body:           requestBody,
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

		// 有界队列 + 批量落库，不再每请求 spawn 一个 goroutine 写库
		enqueueOperationLog(operationLog)
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
