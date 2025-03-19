package system

import (
	"github.com/labstack/echo/v4"
)

type SysUserRouter struct{}

func (s *SysUserRouter) InitSysUserRouter(r *echo.Group) {
	router := r.Group("sys-users")
	{
		router.POST("", sysUserApi.Create) // 创建用户
	}
}
