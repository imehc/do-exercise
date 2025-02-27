package model

import (
	"time"

	"gorm.io/gorm"
)

type IDWrapper struct {
	ID uint `json:"id" gorm:"primarykey;autoIncrement;comment:主键ID"`
}

type ControlWrapper struct {
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	CreatedBy uint           `json:"created_by" gorm:"index;comment:创建者"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	UpdatedBy uint           `json:"updated_by" gorm:"index;comment:更新者"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	DeletedBy uint           `json:"-" gorm:"index;comment:删除者"`
}

// 修改创建者
func (c *ControlWrapper) ModifyCreatedBy(createdBy uint) {
	c.CreatedBy = createdBy
}

// 修改更新者
func (c *ControlWrapper) ModifyUpdatedBy(updatedBy uint) {
	c.UpdatedBy = updatedBy
}

// 修改删除者
func (c *ControlWrapper) ModifyDeletedBy(deletedBy uint) {
	c.DeletedBy = deletedBy
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
