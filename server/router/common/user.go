package common

import (
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (s *AuthRouter) InitUserRouter(r *gin.RouterGroup) gin.IRoutes {
	router := r.Group("user")
	{
		router.GET("bind_email_code", userApi.SendBindEmailCode)     // 发送绑定邮箱验证码
		router.GET("rebind_email_code", userApi.SendRebindEmailCode) // 发送换绑邮箱验证码
		router.PATCH("bind_email", userApi.BindEmail)                // 绑定邮箱
		router.PATCH("rebind_email", userApi.RebindEmail)            // 换绑邮箱
		router.PATCH("modify_password", userApi.UpdatePassword)      // 修改密码
		router.PUT("profile", userApi.UpdateProfile)                 // 修改基本信息
		router.GET("profile", userApi.GetProfile)                    // 获取基本信息
		router.GET("menu", userApi.GetMenu)                          // 获取菜单
	}
	return router
}
