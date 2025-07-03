package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"go.uber.org/zap"
)

type SysInfoApi struct{}

// Get 系统信息
func (s *SysInfoApi) Get(ctx *gin.Context) {
	info, err := sysInfoService.Get()
	if err != nil {
		global.Log.Error("获取系统状态失败", zap.Error(err))
		response.BadRequest(ctx, "getServerStatusFailed")
		return
	}

	response.Success(ctx, info)
}
