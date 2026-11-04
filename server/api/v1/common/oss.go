package common

import (
	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/storage"
	"github.com/imehc/do-exercise/server/util"
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
	data, err := ossService.GetPresignedUrl(util.DB(ctx), userId, tenantId, req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

func (o *OssApi) GetDownloadUrl(ctx *gin.Context) {
	var req common.OssDownloadReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	data, err := ossService.GetDownloadUrl(util.DB(ctx), req.ObjectKey)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

func (o *OssApi) CompleteUpload(ctx *gin.Context) {
	var req common.OssObjectReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	file, err := storage.Complete(util.DB(ctx), ctx.GetString("userId"), req.ObjectKey)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	data := common.OssFileRes{ObjectKey: file.ObjectKey, Size: file.Size, State: file.State}
	if file.Public {
		data.GetObjectUrl = "/" + global.Config.Oss.BucketName + "/" + file.ObjectKey
	}
	response.Success(ctx, data)
}
func (o *OssApi) DeleteFile(ctx *gin.Context) {
	var req common.OssObjectReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := storage.Delete(util.DB(ctx), ctx.GetString("userId"), req.ObjectKey); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.NoContent(ctx)
}
func (o *OssApi) StorageUsage(ctx *gin.Context) {
	charged, limit, err := storage.Usage(util.DB(ctx), ctx.GetString("tenantId"))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, common.OssUsageRes{ChargedBytes: charged, MaxStorageBytes: limit})
}
