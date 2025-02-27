package system

import (
	"github.com/gin-gonic/gin"
)

type DicttRouter struct{}

func (s *AuthRouter) InitDictRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("")

	r1 := r.Group("/dict")
	r2 := r.Group("/dict-data")
	{
		r1.POST("", dictApi.CreateDict)
		r1.DELETE(":dictId", dictApi.DeleteDict)
		r1.PUT(":dictId", dictApi.UpdateDict)
		r1.GET(":dictId", dictApi.GetDict)
		r1.GET("list", dictApi.GetDictList)
	}
	{
		r2.POST("", dictApi.CreateDictData)
		r2.DELETE(":dictDataId", dictApi.DeleteDictData)
		r2.PUT(":dictDataId", dictApi.UpdateDictData)
		r2.GET(":dictDataId", dictApi.GetDictData)
		r2.GET("list", dictApi.GetDictDataList)
	}
	return r
}
