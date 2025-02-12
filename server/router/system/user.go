package system

import (
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *AuthRouter) InitUserRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/user")

	{
		r.POST("", userApi.CreateUser)
		r.DELETE(":userId", userApi.DeleteUser)
		r.PUT(":userId", userApi.UpdateUser)
		r.GET(":userId", userApi.GetUser)
		r.GET("list", userApi.GetUserList)
	}
	return r
}
