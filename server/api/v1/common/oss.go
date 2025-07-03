package common

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
)

type OssApi struct{}

// GetPresignedUrl 获取预签名上传
func (o *OssApi) GetPresignedUrl(ctx *gin.Context) {
	req := common.OssReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := ossService.GetPresignedUrl(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
