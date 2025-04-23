package system

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
)

type SysUser struct {
	Id       int64     `json:"id" gorm:"primarykey;comment:主键"`
	Username string    `json:"username" gorm:"not null;unique;comment:用户名"`
	Nickname string    `json:"nickname" gorm:"comment:昵称"`
	Email    string    `json:"email" gorm:"unique;comment:邮箱"`
	Avatar   string    `json:"avatar" gorm:"comment:头像"`
	Password string    `json:"-" gorm:"not null;comment:密码"`
	Roles    []SysRole `json:"roles" gorm:"many2many:sys_user_role;comment:用户角色"`
	model.ControlWrapper
}

func (*SysUser) TableName() string {
	return "sys_user"
}

func (u *SysUser) BeforeCreate(tx *gorm.DB) (err error) {
	hash := util.Hash{Value: u.Password}
	password, err := hash.Hash()
	if err != nil {
		return err
	}
	u.Password = password

	ctx := tx.Statement.Context
	userId := ctx.Value("userId").(int64)
	u.CreatedBy = userId
	u.UpdatedBy = userId

	return nil
}

func (u *SysUser) BeforeUpdate(tx *gorm.DB) (err error) {
	userId := tx.Statement.Context.Value("userId").(int64)

	if u.UpdatedBy != userId && u.Id != 0 {
		u.UpdatedBy = userId
		err = tx.
			Model(u).
			Select("UpdatedBy").
			Updates(u).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *SysUser) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Id != 0 {
		userId := tx.Statement.Context.Value("userId").(int64)
		u.DeletedBy = userId
		err = tx.
			Model(u).
			Select("DeletedBy").
			Updates(u).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}
