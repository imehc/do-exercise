package system

import (
	"github.com/gin-gonic/gin"
)

type SysRoleRouter struct{}

func (s *SysUserRouter) InitSysRoleRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("roles")
	{
		router.POST("", sysRoleApi.Create)      // 创建角色
		router.DELETE(":id", sysRoleApi.Delete) // 删除角色
		router.PUT(":id", sysRoleApi.Update)    // 更新角色
		router.GET(":id", sysRoleApi.Get)       // 获取单个角色
		router.GET("", sysRoleApi.GetList)      // 获取角色列表
	}
	return router
}
