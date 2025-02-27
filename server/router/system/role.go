package system

import (
	"github.com/gin-gonic/gin"
)

type RoleRouter struct{}

func (s *AuthRouter) InitRoleRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/roles")

	{
		r.GET("", roleApi.GetRoleList)
		r.POST("", roleApi.CreateRole)
		r.DELETE(":roleId", roleApi.DeleteRole)
		r.PUT(":roleId", roleApi.UpdateRole)
		r.PUT(":roleId/data-scope", roleApi.UpdateRoleDataScope)
		r.PUT(":roleId/menu-scope", roleApi.UpdateMenuDataScope)
		r.PUT(":roleId/api-scope", roleApi.UpdateApiDataScope)
		r.GET(":roleId", roleApi.GetRole)
	}
	return r
}
