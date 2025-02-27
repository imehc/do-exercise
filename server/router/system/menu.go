package system

import (
	"github.com/gin-gonic/gin"
)

type MenuRouter struct{}

func (s *AuthRouter) InitMenuRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/menu")

	{
		r.POST("", menuApi.CreateMenu)
		r.DELETE(":menuId", menuApi.DeleteMenu)
		r.PUT(":menuId", menuApi.UpdateMenu)
		r.GET(":menuId", menuApi.GetMenu)
		r.GET("list", menuApi.GetMenuList)
	}
	return r
}
