package system

import (
	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

type SysRole struct {
	model.IdWrapper

	Name  string    `json:"name" gorm:"not null;comment:角色名称"`
	Code  string    `json:"code" gorm:"not null;unique;comment:角色编码"`
	Menus []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;comment:角色菜单"`

	model.ControlWrapper
}

func (r *SysRole) BeforeCreate(tx *gorm.DB) (err error) {
	userId := tx.Statement.Context.Value("userId").(int64)
	r.CreatedBy = userId
	r.UpdatedBy = userId

	return nil
}

func (r *SysRole) BeforeUpdate(tx *gorm.DB) (err error) {
	userId := tx.Statement.Context.Value("userId").(int64)

	if r.UpdatedBy != userId && r.Id != 0 {
		r.UpdatedBy = userId
		err = tx.
			Model(r).
			Select("UpdatedBy").
			Updates(r).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SysRole) BeforeDelete(tx *gorm.DB) (err error) {
	if r.Id == 0 {
		userId := tx.Statement.Context.Value("userId").(int64)
		r.DeletedBy = userId
		err = tx.
			Model(r).
			Select("DeletedBy").
			Updates(r).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
