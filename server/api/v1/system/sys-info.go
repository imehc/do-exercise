package system

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/util"
)

type SysInfoApi struct{}

// Get 系统信息
func (s *SysInfoApi) Get(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	// 创建一个定时器，每隔指定时间触发一次
	ticker := time.NewTicker(time.Duration(global.Config.System.InfoInterval) * time.Second)
	defer ticker.Stop()

	// 创建一个用于检测客户端断开连接的通道
	done := ctx.Done()

	// 首次立即发送数据
	if info, err := sysInfoService.Get(); !sendSSEMessage(ctx, info, err) {
		return
	}

	// 循环发送数据
	for {
		select {
		case <-done:
			// 客户端断开连接
			return
		case <-ticker.C:
			// 每隔指定时间获取一次新数据
			if info, err := sysInfoService.Get(); !sendSSEMessage(ctx, info, err) {
				return
			}
		}
	}
}

// SendSSEMessage SSE 消息发送函数
func sendSSEMessage(ctx *gin.Context, data *util.SystemInfo, err error) bool {
	lang := ctx.GetString("lang")
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
		return false
	}

	ctx.SSEvent("message", data)
	ctx.Writer.Flush()
	return true
}
