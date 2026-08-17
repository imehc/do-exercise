package system

import (
	"github.com/gin-gonic/gin"
)

type SysTenantRouter struct{}

func (s *SysTenantRouter) InitSysTenantRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("tenants")
	{
		router.GET("", tenantApi.GetList)
		router.POST("", tenantApi.Create)
		// 可被选作租户管理员的现有用户（创建租户时用，需在 :id 之前注册避免被吞）
		router.GET("assignable_admins", tenantApi.ListAssignableAdmins)
		router.GET(":id", tenantApi.Get)
		router.PUT(":id", tenantApi.Update)
		router.DELETE(":id", tenantApi.Delete)
		// 租户分配用户
		router.GET(":id/assignable_users", tenantApi.ListAssignableUsers)
		router.POST(":id/users", tenantApi.AssignUsers)
	}
	return r
}