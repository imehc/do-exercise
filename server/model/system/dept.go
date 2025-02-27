package system

import "github.com/imehc/do-exercise/server/model"

type Dept struct {
	model.IDWrapper

	ParentId  int    `gorm:"comment:上级部门" json:"parent_id"`
	Path      string `gorm:"size:255;" json:"path"`
	Name      string `gorm:"size:128;comment:部门名称" json:"name"`
	Leader    string `gorm:"size:128;comment:负责人" json:"leader"`
	Phone     string `gorm:"size:11;comment:联系电话" json:"phone"`
	Email     string `gorm:"size:64;comment:邮箱" json:"email"`
	DataScope string `gorm:"-" json:"data_scope"`
	Params    string `gorm:"-" json:"params"`
	Children  []Dept `gorm:"-" json:"children"`

	model.SortWrapper
	model.StatusWrapper
	model.ControlWrapper
}

func (*Dept) TableName() string {
	return "sys_dept"
}
