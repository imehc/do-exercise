package system

import (
	"github.com/gin-gonic/gin"
)

type RoleRouter struct{}

func (s *AuthRouter) InitRoleRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/role")

	{
		r.POST("", roleApi.CreateRole)
		r.DELETE(":roleId", roleApi.DeleteRole)
		r.PUT(":roleId", roleApi.UpdateRole)
		r.GET(":roleId", roleApi.GetRole)
		r.GET("list", roleApi.GetRoleList)
	}
	return r
}
