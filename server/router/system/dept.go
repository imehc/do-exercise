package system

import (
	"github.com/gin-gonic/gin"
)

type DeptRouter struct{}

func (s *AuthRouter) InitDeptRouter(router *gin.RouterGroup) gin.IRoutes {
	r := router.Group("/dept")

	{
		r.POST("", deptApi.CreateDept)
		r.DELETE(":deptId", deptApi.DeleteDept)
		r.PUT(":deptId", deptApi.UpdateDept)
		r.GET(":deptId", deptApi.GetDept)
		r.GET("list", deptApi.GetDeptList)
		r.GET("tree", deptApi.GetDeptTree)
	}
	return r
}
