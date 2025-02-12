package system

import (
	"github.com/gin-gonic/gin"
)

type PostRouter struct{}

func (s *AuthRouter) InitPostRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/post")

	{
		r.POST("", postApi.CreatePost)
		r.DELETE(":postId", postApi.DeletePost)
		r.PUT(":postId", postApi.UpdatePost)
		r.GET(":postId", postApi.GetPost)
		r.GET("list", postApi.GetPostList)
	}
	return r
}
