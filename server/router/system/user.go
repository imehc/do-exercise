package system

import (
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *AuthRouter) InitUserRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/users")

	{
		r.GET("", userApi.GetUserList)
		r.POST("", userApi.CreateUser)
		r.DELETE(":userId", userApi.DeleteUser)
		r.PUT(":userId", userApi.UpdateUser)
		r.GET(":userId", userApi.GetUser)
	}
	return r
}

func (s *AuthRouter) InitCurrentUserRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/user")

	{
		r.POST("", userApi.UpdateUser)
		r.GET("", userApi.GetUserInfo)
	}
	return r
}
