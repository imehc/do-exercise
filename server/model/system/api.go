package system

import "github.com/imehc/do-exercise/server/model"

type SysApi struct {
	model.IDWrapper

	Handle string `gorm:"size:128;comment:handle" json:"handle"`
	Title  string `gorm:"size:128;comment:标题" json:"title"`
	Path   string `gorm:"size:128;comment:地址" json:"path"`
	Type   string `gorm:"size:16;comment:接口类型" json:"type"`
	Action string `gorm:"size:16;comment:请求类型" json:"action"`

	model.ControlWrapper
}

func (SysApi) TableName() string {
	return "sys_api"
}
