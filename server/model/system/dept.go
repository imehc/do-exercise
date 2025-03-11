package system

import "github.com/imehc/do-exercise/server/model"

type Dept struct {
	DeptId   uint   `gorm:"primarykey;autoIncrement;comment:主键ID" json:"id"`
	ParentId *uint  `gorm:"comment:上级部门" json:"parent_id"`
	Path     string `gorm:"size:255;" json:"path"`
	Name     string `gorm:"size:128;comment:部门名称" json:"name"`
	Posts    []Post `gorm:"many2many:sys_dept_post;foreignKey:DeptId;joinForeignKey:dept_id;references:PostId;joinReferences:post_id" json:"posts"`
	PostIds  []int  `gorm:"-" json:"post_ids"`
	Children []Dept `gorm:"-" json:"children"`

	model.SortWrapper
	model.StatusWrapper
	model.ControlWrapper
}

func (*Dept) TableName() string {
	return "sys_dept"
}
