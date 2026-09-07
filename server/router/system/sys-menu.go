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
	}
	return router
}

// InitSysMenuTreeRouter 注册菜单树：认证即可读，不经 Casbin 鉴权。
//
// 菜单树不是「菜单管理」这个页面的私有数据，而是整套权限模型的字典：角色授权弹窗
// 要靠它勾选权限。若继续由「菜单管理-查询」按钮授权，「能管角色但没有菜单查询权」
// 的角色一打开授权弹窗就是 403 白屏（P2-2）——两个互不相干的权限被隐式耦合了。
//
// 放开的是「平台统一维护的菜单目录」本身（名称/路由/权限标识），里面没有任何租户
// 业务数据；可见范围仍按 scope 与当前上下文过滤。菜单的增删改仍在 Casbin 组里。
func (s *SysMenuRouter) InitSysMenuTreeRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("menus")
	{
		router.GET("tree", sysMenuApi.GetTree) // 获取菜单树
	}
	return router
}
