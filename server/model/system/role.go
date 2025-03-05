package system

import "github.com/imehc/do-exercise/server/model"

type Role struct {
	RoleId    uint   `gorm:"primarykey;autoIncrement;comment:主键ID" json:"role_id"`
	Name      string `gorm:"size:128;not null;comment:角色名称" json:"name"`
	Key       string `gorm:"size:128;comment:角色权限字符串" json:"key"`
	IsAdmin   bool   `gorm:"comment:是否为admin角色" json:"is_admin"`
	DataScope uint   `gorm:"size:128;default:1;comment:数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）" json:"data_scope"`
	MenuIds   []int  `gorm:"-" json:"menu_ids"` // 菜单权限
	ApiIds    []int  `gorm:"-" json:"api_ids"`  // 接口权限
	Menus     []Menu `gorm:"many2many:sys_role_menu;foreignKey:RoleId;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id" json:"menus"`
	Apis      []Api  `gorm:"many2many:sys_role_api;foreignKey:RoleId;joinForeignKey:role_id;references:ApiId;joinReferences:api_id" json:"apis"`

	model.SortWrapper
	model.RemarkWrapper
	model.StatusWrapper
	model.ControlWrapper
}

func (*Role) TableName() string {
	return "sys_role"
}
