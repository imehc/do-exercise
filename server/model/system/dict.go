package system

import (
	"github.com/imehc/do-exercise/server/model"
)

type Dict struct {
	model.IDWrapper

	Name   string `gorm:"not null;size:128;comment:名称" json:"name"`             // 字典名称
	Type   string `gorm:"uniqueIndex;not null;size:128;comment:类型" json:"type"` // 字典类型
	Remark string `gorm:"size:255;comment:备注;" json:"remark"`                   // 备注

	model.StatusWrapper
	model.ControlWrapper
}

func (Dict) TableName() string {
	return "sys_dict"
}
