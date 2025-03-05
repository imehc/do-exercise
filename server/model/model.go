package model

import (
	"time"

	"gorm.io/gorm"
)

type ControlWrapper struct {
	CreateAt time.Time      `json:"create_at" gorm:"comment:创建时间"`
	CreateBy uint           `json:"create_by" gorm:"index;comment:创建者"`
	UpdateAt time.Time      `json:"update_at" gorm:"comment:更新时间"`
	UpdateBy uint           `json:"update_by" gorm:"index;comment:更新者"`
	DeleteAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	DeleteBy uint           `json:"-" gorm:"index;comment:删除者"`
}

// 修改创建者
func (c *ControlWrapper) ModifyCreateBy(createBy uint) {
	c.CreateBy = createBy
}

// 修改更新者
func (c *ControlWrapper) ModifyUpdateBy(updateBy uint) {
	c.UpdateBy = updateBy
}

// 修改删除者
func (c *ControlWrapper) ModifyDeleteBy(deleteBy uint) {
	c.DeleteBy = deleteBy
}

type StatusWrapper struct {
	Status int `gorm:"default:1;size:4;comment:状态" json:"status"` // 0未启用 1启用 2禁用
}

type SortWrapper struct {
	Sort int `gorm:"size:4;comment:排序" json:"sort"`
}

type RemarkWrapper struct {
	Remark string `gorm:"size:255;comment:备注;" json:"remark"`
}
