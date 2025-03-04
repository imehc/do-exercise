package system

import "github.com/imehc/do-exercise/server/model"

type Role struct {
	model.IDWrapper

	Name      string `gorm:"size:128;not null;comment:角色名称" json:"name"`
	Key       string `gorm:"size:128;comment:角色权限字符串" json:"key"`
	Flag      string `gorm:"size:128;" json:"flag"`
	IsAdmin   bool   `gorm:"comment:是否为admin角色" json:"is_admin"`
	DataScope uint   `gorm:"size:128;default:1;comment:数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）" json:"data_scope"`
	Params    string `gorm:"-" json:"params"`
	MenuIds   []int  `gorm:"-" json:"menu_ids"` // 菜单权限
	DeptIds   []int  `gorm:"-" json:"dept_ids"` // 部门权限
	ApiIds    []int  `gorm:"-" json:"api_ids"`  // 接口权限
	Menus     []Menu `gorm:"many2many:sys_role_menu;foreignKey:ID;joinForeignKey:role_id;references:ID;joinReferences:menu_id" json:"menus"`
	Depts     []Dept `gorm:"many2many:sys_role_dept;foreignKey:ID;joinForeignKey:role_id;references:ID;joinReferences:dept_id" json:"depts"`
	Apis      []Api  `gorm:"many2many:sys_role_api;foreignKey:ID;joinForeignKey:role_id;references:ID;joinReferences:api_id" json:"apis"`

	model.SortWrapper
	model.RemarkWrapper
	model.StatusWrapper
	model.ControlWrapper
}

func (*Role) TableName() string {
	return "sys_role"
}
