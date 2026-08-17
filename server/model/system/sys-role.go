package system

import (
	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

type SysRole struct {
	model.IdWrapper

	Name     string    `json:"name" gorm:"not null;comment:角色名称"`
	Code     string    `json:"code" gorm:"not null;uniqueIndex:idx_sys_role_code_tenant,priority:1;comment:角色编码"`
	TenantId string    `json:"tenant_id" gorm:"column:tenant_id;type:varchar(32);not null;default:'';index;comment:租户ID;uniqueIndex:idx_sys_role_code_tenant,priority:2"`
	Menus    []SysMenu `json:"menus" gorm:"many2many:sys_role_menu;comment:角色菜单"`

	model.ControlWrapper
}

func (r *SysRole) BeforeCreate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)
	r.CreatedBy = userId
	r.UpdatedBy = userId

	return nil
}

func (r *SysRole) BeforeUpdate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)
	if userId == "" {
		return nil
	}

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
		return nil
	}
	userId := model.CurrentUserID(tx)
	r.DeletedBy = userId
	err = tx.
		Model(r).
		Select("DeletedBy").
		Updates(r).
		Error
	if err != nil {
		return err
	}

	return nil
}
