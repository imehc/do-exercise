package common

import (
	"github.com/gin-gonic/gin"
)

type OssRouter struct{}

func (s *AuthRouter) InitOssRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("")
	{
		router.GET("get_presigned_url", ossApi.GetPresignedUrl)
		router.POST("files/complete", ossApi.CompleteUpload)
		router.DELETE("files", ossApi.DeleteFile)
		router.GET("files/usage", ossApi.StorageUsage)
		router.GET("get_download_url", ossApi.GetDownloadUrl) // 获取预上传签名
	}
	return router
}
