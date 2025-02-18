package system

import "github.com/imehc/do-exercise/server/model"

type Menu struct {
	model.IDWrapper

	Name       string   `gorm:"size:128;comment:菜单名称" json:"name"`
	Title      string   `gorm:"size:128;comment:菜单标题" json:"title"`
	Icon       string   `gorm:"size:128;comment:菜单图标" json:"icon"`
	Path       string   `gorm:"size:128;comment:路由地址" json:"path"`
	Paths      string   `gorm:"size:128;comment:组件路径" json:"paths"`
	Type       string   `gorm:"size:1;comment:菜单类型（M目录 C菜单 F按钮）" json:"type"`
	Action     string   `gorm:"size:16;comment:请求方法" json:"action"`
	Permission string   `gorm:"size:255;comment:权限标识" json:"perms"`
	ParentId   int      `gorm:"comment:父菜单ID" json:"parent_id"`
	NoCache    bool     `gorm:"comment:是否缓存" json:"no_cache"`
	Component  string   `gorm:"size:128;comment:组件路径" json:"component"`
	Visible    int      `gorm:"size:1;default:0;comment:显示状态（0显示 1隐藏）" json:"visible"`
	IsFrame    int      `gorm:"default:1;comment:是否为外链（0是 1否）" json:"is_frame"`
	IsCache    int      `gorm:"default:0;comment:是否缓存（0缓存 1不缓存）" json:"is_cache"`
	SysApi     []SysApi `gorm:"many2many:sys_menu_api_rule" json:"sys_api"`
	Apis       []int    `gorm:"-" json:"apis"`
	DataScope  string   `gorm:"-" json:"data_scope"`
	Params     string   `gorm:"-" json:"params"`
	RoleId     int      `gorm:"-" json:"role_id"`
	Children   []Menu   `gorm:"-" json:"children"`
	IsSelect   bool     `gorm:"-" json:"is_select"`

	model.SortWrapper
	model.ControlWrapper
}

func (*Menu) TableName() string {
	return "sys_menu"
}
