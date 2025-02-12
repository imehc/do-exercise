package system

import (
	"github.com/gin-gonic/gin"
)

type DicttRouter struct{}

func (s *AuthRouter) InitDictRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/dict")

	{
		r.POST("", dictApi.CreateDict)
		r.DELETE(":dictId", dictApi.DeleteDict)
		r.PUT(":dictId", dictApi.UpdateDict)
		r.GET(":dictId", dictApi.GetDict)
		r.GET("list", dictApi.GetDictList)
	}
	{
		r.POST(":dictId", dictApi.CreateDictDetail)
		r.DELETE(":dictId/:dictDetailId", dictApi.DeleteDictDetail)
		r.PUT(":dictId/:dictDetailId", dictApi.UpdateDictDetail)
		r.GET(":dictId/:dictDetailId", dictApi.GetDictDetail)
		r.GET(":dictId/list", dictApi.GetDictDetailList)
	}
	return r
}
