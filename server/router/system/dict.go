package system

import (
	"github.com/gin-gonic/gin"
)

type DicttRouter struct{}

func (s *AuthRouter) InitDictRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("")

	r1 := r.Group("/dicts")
	r2 := r.Group("/dict-data")
	{
		r1.GET("", dictApi.GetDictList)
		r1.POST("", dictApi.CreateDict)
		r1.DELETE(":dictId", dictApi.DeleteDict)
		r1.PUT(":dictId", dictApi.UpdateDict)
		r1.GET(":dictId", dictApi.GetDict)
	}
	{
		r2.GET("", dictApi.GetDictDataList)
		r2.POST("", dictApi.CreateDictData)
		r2.DELETE(":dictDataId", dictApi.DeleteDictData)
		r2.PUT(":dictDataId", dictApi.UpdateDictData)
		r2.GET(":dictDataId", dictApi.GetDictData)
	}
	return r
}
