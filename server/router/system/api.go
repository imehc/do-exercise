package system

import (
	"github.com/gin-gonic/gin"
)

type ApiRouter struct{}

func (s *AuthRouter) InitApiRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/apis")

	{
		r.GET("", apiApi.GetApiList)
		r.POST("", apiApi.CreateApi)
		r.DELETE(":apiId", apiApi.DeleteApi)
		r.PUT(":apiId", apiApi.UpdateApi)
		r.GET(":apiId", apiApi.GetApi)
		r.GET("list", apiApi.GetApis)
	}
	return r
}
