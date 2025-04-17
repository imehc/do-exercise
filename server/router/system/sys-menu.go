package system

import "github.com/gin-gonic/gin"

type SysMenuRouter struct{}

func (s *SysMenuRouter) InitSysMenuRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("menus")
	{
		router.POST("", sysMenuApi.Create)      // 创建菜单
		router.DELETE(":id", sysMenuApi.Delete) // 删除菜单
		router.PUT(":id", sysMenuApi.Update)    // 更新菜单
		router.GET(":id", sysMenuApi.Get)       // 获取单个菜单
		router.GET("tree", sysMenuApi.GetTree)  // 获取菜单树
	}
	return router
}
