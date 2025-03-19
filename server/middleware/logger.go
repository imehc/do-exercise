package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var Logger = middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	LogURI:       true,
	LogStatus:    true,
	LogLatency:   true,
	LogMethod:    true,
	LogRequestID: true,
	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
		zap.L().Info("HTTP Request",
			zap.String("method", v.Method),
			zap.String("URI", v.URI),
			zap.Int("status", v.Status),
			zap.Duration("latency", v.Latency),
			zap.String("request_id", v.RequestID),
			zap.String("remote_ip", c.RealIP()),
			zap.String("host", c.Request().Host),
			zap.String("user_agent", c.Request().UserAgent()),
		)

		// 对于错误状态码，额外记录错误日志
		if v.Status >= 400 {
			// 获取错误信息
			var errMsg string
			if err := c.Get("error"); err != nil {
				if e, ok := err.(error); ok {
					errMsg = e.Error()
				} else if s, ok := err.(string); ok {
					errMsg = s
				} else {
					errMsg = "未知错误"
				}
			} else {
				errMsg = "未知错误"
			}

			zap.L().Error("HTTP Error",
				zap.String("method", v.Method),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
				zap.String("request_id", v.RequestID),
				zap.String("error", errMsg),
			)
		}

		return nil
	},
})
