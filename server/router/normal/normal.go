package normal

import "github.com/gin-gonic/gin"

type NormalRouter struct{}

func (s *NormalRouter) InitSysUserRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("")
	{
	}
	return router
}
