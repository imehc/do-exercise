package model

import (
	"time"

	"gorm.io/gorm"
)

type IdWrapper struct {
	Id uint `json:"id" gorm:"primarykey;comment:主键"`
}

type ControlWrapper struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by;type:varchar(32);comment:创建人"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by;type:varchar(32);comment:更新人"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index;comment:删除时间"`
	DeletedBy string         `json:"deleted_by" gorm:"column:deleted_by;type:varchar(32);comment:删除人"`
}
