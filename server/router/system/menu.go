package system

import (
	"github.com/gin-gonic/gin"
)

type MenuRouter struct{}

func (s *AuthRouter) InitMenuRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/menus")

	{
		r.GET("", menuApi.GetMenuList)
		r.POST("", menuApi.CreateMenu)
		r.DELETE(":menuId", menuApi.DeleteMenu)
		r.PUT(":menuId", menuApi.UpdateMenu)
		r.GET(":menuId", menuApi.GetMenu)
		r.GET("tree", menuApi.GetMenuTree)
		r.GET("compact", menuApi.GetMenuTreeCompact)
	}
	return r
}
