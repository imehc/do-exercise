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
	// 系统信息是宿主机层面的运维数据（CPU/内存/磁盘），无法按租户切分，
	// 因此只对平台超级管理员开放。菜单树里也已按 scope=platform 隐藏，
	// 这里再挡一次，避免旧租户遗留的 Casbin 策略仍能直接调到该端点。
	if !ctx.GetBool("isSuperAdmin") {
		response.Forbidden(ctx, global.I18.Translate("platformOnly", ctx.GetString("lang")))
		return
	}

	info, err := sysInfoService.Get()
	if err != nil {
		global.Log.Error("获取系统状态失败", zap.Error(err))
		response.BadRequest(ctx, "getServerStatusFailed")
		return
	}

	response.Success(ctx, info)
}
