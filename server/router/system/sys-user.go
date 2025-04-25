package system

import (
	"github.com/gin-gonic/gin"
)

type SysUserRouter struct{}

func (s *SysUserRouter) InitSysUserRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("users")
	{
		router.POST("", sysUserApi.Create)                         // 创建用户
		router.DELETE(":id", sysUserApi.Delete)                    // 删除用户
		router.PUT(":id", sysUserApi.Update)                       // 更新用户
		router.GET(":id", sysUserApi.Get)                          // 获取单个用户
		router.GET("", sysUserApi.GetList)                         // 获取用户列表
		router.PUT(":id/reset_password", sysUserApi.ResetPassword) // 重置用户密码
	}
	return router
}
