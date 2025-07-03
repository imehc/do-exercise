package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type SysTokenApi struct{}

// GetList 获取token列表
func (s *SysTokenApi) FindAll(ctx *gin.Context) {
	data, err := sysTokenService.FindAll()
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

// Delete 删除token
func (s *SysTokenApi) Delete(ctx *gin.Context) {
	var req request.SysTokenDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := sysTokenService.Delete(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}

// ModityStatus 修改token状态
func (s *SysTokenApi) ModityStatus(ctx *gin.Context) {
	var req request.SysTokenModityStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := sysTokenService.ModityStatus(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}
