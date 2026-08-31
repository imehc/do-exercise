package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/util"
)

type SysOperationLogApi struct{}

// GetList 获取操作日志列表。
// 传请求级 DB，让租户插件把可见范围限制在当前租户；平台超级管理员看全量。
func (s *SysOperationLogApi) GetList(ctx *gin.Context) {
	var req common.Pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := sysOperationLogService.GetList(util.DB(ctx), req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
