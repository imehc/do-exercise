package common

import (
	"github.com/gin-gonic/gin"
)

type OssRouter struct{}

func (s *AuthRouter) InitOssRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("")
	{
		router.GET("get_presigned_url", ossApi.GetPresignedUrl) // 获取预上传签名
	}
	return router
}
