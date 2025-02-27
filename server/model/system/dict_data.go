package system

import "github.com/imehc/do-exercise/server/model"

type DictData struct {
	model.IDWrapper
	DictType string `gorm:"not null;size:64;comment:字典类型" json:"dict_type"`
	Label    string `gorm:"not null;size:128;comment:数据标签" json:"label"`
	Value    string `gorm:"not null;size:255;comment:数据键值" json:"value"`

	model.RemarkWrapper
	model.SortWrapper
	model.StatusWrapper
	model.ControlWrapper
}

func (DictData) TableName() string {
	return "sys_dict_data"
}
