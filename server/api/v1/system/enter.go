package system

import "github.com/imehc/do-exercise/server/service"

type ApiGroup struct {
	AuthApi
	CaptchaApi
	DeptApi
	DictApi
	MenuApi
	PostApi
	RoleApi
	UserApi
	ApiApi
}

var (
	authService = service.ServiceGroupApp.SystemServiceGroup.AuthService // 登录服务
	deptService = service.ServiceGroupApp.SystemServiceGroup.DeptService // 部门服务
	dictService = service.ServiceGroupApp.SystemServiceGroup.DictService // 字典服务
	menuService = service.ServiceGroupApp.SystemServiceGroup.MenuService // 菜单服务
	postService = service.ServiceGroupApp.SystemServiceGroup.PostService // 岗位服务
	roleService = service.ServiceGroupApp.SystemServiceGroup.RoleService // 角色服务
	userService = service.ServiceGroupApp.SystemServiceGroup.UserService // 用户服务
	apiService  = service.ServiceGroupApp.SystemServiceGroup.ApiService  // 用户服务
)
