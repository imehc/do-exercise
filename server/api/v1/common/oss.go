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
	// 对象 key 由服务端按用户生成，客户端不能指定，避免覆盖他人对象
	userId := ctx.MustGet("userId").(string)
	tenantId := ctx.GetString("tenantId")
	data, err := ossService.GetPresignedUrl(userId, tenantId, req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}
