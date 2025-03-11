package system

import "github.com/imehc/do-exercise/server/model"

type Menu struct {
	MenuId     uint   `gorm:"primarykey;autoIncrement;comment:主键ID" json:"id"`
	Name       string `gorm:"size:128;comment:路由名称" json:"name"`
	Title      string `gorm:"size:128;comment:菜单标题" json:"title"`
	Icon       string `gorm:"size:128;comment:菜单图标" json:"icon"`
	Route      string `gorm:"size:128;comment:路由地址" json:"route"`
	Path       string `gorm:"size:128;comment:组件路径" json:"path"`
	Type       string `gorm:"size:1;comment:菜单类型（M目录 C菜单 F按钮）" json:"type"`
	Action     string `gorm:"size:16;comment:请求方法" json:"action"`
	Permission string `gorm:"size:255;comment:权限标识" json:"perms"`
	ParentId   uint   `gorm:"comment:父菜单ID" json:"parent_id"`
	Visible    bool   `gorm:"default:false;comment:是否可见" json:"visible"`
	Apis       []Api  `gorm:"many2many:sys_menu_api;foreignKey:MenuId;joinForeignKey:menu_id;references:ApiId;joinReferences:api_id" json:"apis"`
	ApiIds     []int  `gorm:"-" json:"api_ids"`
	Children   []Menu `gorm:"-" json:"children"`

	model.SortWrapper
	model.ControlWrapper
}

func (*Menu) TableName() string {
	return "sys_menu"
}
