package system

import "github.com/imehc/do-exercise/server/model"

type SysRole struct {
	model.IdWrapper

	Name  string    `json:"name" gorm:"not null;comment:角色名称"`
	Code  string    `json:"code" gorm:"not null;unique;comment:角色编码"`
	Menus []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;comment:角色菜单"`

	model.ControlWrapper
}
