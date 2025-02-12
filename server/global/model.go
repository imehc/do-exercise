package global

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `json:"id" gorm:"primarykey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

type ControlBy struct {
	CreatedBy uint `json:"created_by" gorm:"index;comment:创建者"`
	UpdatedBy uint `json:"updated_by" gorm:"index;comment:更新者"`
	DeletedBy uint `json:"deleted_by" gorm:"index;comment:删除者"`
}

// 修改创建者
func (c *ControlBy) ModifyCreatedBy(createdBy uint) {
	c.CreatedBy = createdBy
}

// 修改更新者
func (c *ControlBy) ModifyUpdatedBy(updatedBy uint) {
	c.UpdatedBy = updatedBy
}

// 修改删除者
func (c *ControlBy) ModifyDeletedBy(deletedBy uint) {
	c.DeletedBy = deletedBy
}
