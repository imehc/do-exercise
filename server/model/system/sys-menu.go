package system

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

type SysMenu struct {
	model.IdWrapper

	Name       string   `json:"name" gorm:"not null;comment:菜单名称"`
	ParentId   *uint    `json:"parent_id" gorm:"default:null;comment:父菜单ID"`
	Permission *string  `json:"permission" gorm:"size:64;default:null;unique;comment:权限标识"`
	Icon       string   `json:"icon" gorm:"size:64;default:null;comment:菜单图标"`
	Type       uint8    `json:"type" gorm:"not null;comment:菜单类型(1:目录,2:菜单,3:按钮)"`
	Route      string   `json:"route" gorm:"size:128;default:null;comment:菜单路由"`
	Component  string   `json:"component" gorm:"size:128;default:null;comment:组件地址"`
	Sort       uint     `json:"sort" gorm:"default:0;comment:显示顺序"`
	Visible    bool     `json:"visible" gorm:"default:false;comment:是否显示"`
	Apis       []SysApi `json:"apis" gorm:"many2many:sys_menu_apis;"`

	model.ControlWrapper
}

func (m *SysMenu) BeforeCreate(tx *gorm.DB) (err error) {
	userId := tx.Statement.Context.Value(global.ContextUserIDKey).(int64)
	m.CreatedBy = userId
	m.UpdatedBy = userId

	return nil
}

func (m *SysMenu) BeforeUpdate(tx *gorm.DB) (err error) {
	userId := tx.Statement.Context.Value(global.ContextUserIDKey).(int64)

	if m.UpdatedBy != userId && m.Id != 0 {
		m.UpdatedBy = userId
		err = tx.
			Model(m).
			Select("UpdatedBy").
			Updates(m).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SysMenu) BeforeDelete(tx *gorm.DB) (err error) {
	if m.Id == 0 {
		userId := tx.Statement.Context.Value(global.ContextUserIDKey).(int64)
		m.DeletedBy = userId
		err = tx.
			Model(m).
			Select("DeletedBy").
			Updates(m).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
