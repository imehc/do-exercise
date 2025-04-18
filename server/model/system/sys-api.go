package system

import (
	"time"

	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

type SysApi struct {
	model.IdWrapper
	Path        string         `json:"path" gorm:"not null;comment:api路径"`
	Description string         `json:"description" gorm:"size:255;default:null;comment:描述"`
	Group       string         `json:"group" gorm:"default:null;index;comment:api组"`
	Method      string         `json:"method" gorm:"not null;comment:方法"`
	Disabled    bool           `json:"disabled" gorm:"default:false;comment:是否禁用"`
	Sort        uint           `json:"sort" gorm:"default:0;comment:排序"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index;comment:删除时间"`
}
