package system

import (
	"github.com/gin-gonic/gin"
)

type SysUserRouter struct{}

func (s *SysUserRouter) InitSysUserRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("users")
	{
		router.GET("check_username", sysUserApi.UsernameExists)                // 检查用户名是否已存在
		router.POST("", sysUserApi.Create)                                     // 创建用户
		router.DELETE(":id", sysUserApi.Delete)                                // 删除用户
		router.PUT(":id", sysUserApi.Update)                                   // 更新用户
		router.GET(":id", sysUserApi.Get)                                      // 获取单个用户
		router.GET("", sysUserApi.GetList)                                     // 获取用户列表
		router.PUT(":id/reset_password", sysUserApi.ResetPassword)             // 重置用户密码
		router.PATCH(":id/tenant", sysUserApi.AssignTenant)                    // 分配用户到租户（仅超管）
		router.GET(":id/assignable_tenants", sysUserApi.ListAssignableTenants) // 可分配候选租户（仅超管）
	}
	return router
}
