package system

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/common/status"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type SysTokenApi struct{}

// GetList 获取token列表
func (s *SysTokenApi) FindAll(ctx *gin.Context) {
	lang := ctx.GetString("lang")

	data, err := sysTokenService.FindAll()
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
	}
	response.Success(ctx, data)
}

// Delete 删除token
func (s *SysTokenApi) Delete(ctx *gin.Context) {
	lang := ctx.GetString("lang")

	var req request.SysTokenDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := sysTokenService.Delete(req)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
	}
	response.NoContent(ctx)
}

// ModityStatus 修改token状态
func (s *SysTokenApi) ModityStatus(ctx *gin.Context) {
	lang := ctx.GetString("lang")

	var req request.SysTokenModityStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	err := sysTokenService.ModityStatus(req)
	if err != nil {
		response.BadRequest(ctx, response.ValidationError{
			Type:    status.BAD_REQUEST_MSG,
			Message: global.I18.Translate(err.Error(), lang),
		})
	}
	response.NoContent(ctx)
}
