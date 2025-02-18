package system

import "github.com/imehc/do-exercise/server/model"

type Post struct {
	model.IDWrapper

	Name   string `gorm:"size:128;not null;comment:岗位名称" json:"name"`
	Code   string `gorm:"size:128;not null;comment:岗位编码" json:"code"`
	DeptId int    `gorm:"comment:所属部门" json:"dept_id"`
	Dept   Dept   `gorm:"foreignKey:DeptId" json:"dept"`

	model.SortWrapper
	model.StatusWrapper
	model.RemarkWrapper
	model.ControlWrapper
}

func (*Post) TableName() string {
	return "sys_post"
}
