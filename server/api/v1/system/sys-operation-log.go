package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
)

type SysOperationLogApi struct{}

// GetList 获取操作日志列表
func (s *SysOperationLogApi) GetList(ctx *gin.Context) {
	var req common.Pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := sysOperationLogService.GetList(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
