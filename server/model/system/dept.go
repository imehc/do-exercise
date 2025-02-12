package system

import "github.com/imehc/do-exercise/server/global"

type Dept struct {
	global.Model
	global.ControlBy

	ParentId  *int   `json:"parent_id" gorm:"comment:上级部门"`
	Path      string `json:"path" gorm:"size:255;"`
	Name      string `json:"name"  gorm:"size:128;comment:部门名称"`
	Sort      int    `json:"sort" gorm:"size:4;comment:排序"`
	Leader    string `json:"leader" gorm:"size:128;comment:负责人"`
	Phone     string `json:"phone" gorm:"size:11;comment:联系电话"`
	Email     string `json:"email" gorm:"size:64;comment:邮箱"`
	Status    int    `json:"status" gorm:"size:4;comment:状态"`
	DataScope string `json:"data_scope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
	Children  []Dept `json:"children" gorm:"-"`
}

func (*Dept) TableName() string {
	return "sys_dept"
}
