package system

import "github.com/imehc/do-exercise/server/model"

type Post struct {
	PostId uint   `gorm:"primarykey;autoIncrement;comment:主键ID" json:"post_id"`
	Name   string `gorm:"size:128;not null;comment:岗位名称" json:"name"`
	Code   string `gorm:"size:128;not null;comment:岗位编码" json:"code"`

	model.SortWrapper
	model.StatusWrapper
	model.RemarkWrapper
	model.ControlWrapper
}

func (*Post) TableName() string {
	return "sys_post"
}
