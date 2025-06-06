package system

import (
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
)

type SysUser struct {
	UserId   string    `json:"id" gorm:"primarykey;column:id;type:varchar(32);comment:主键"`
	Username string    `json:"username" gorm:"not null;unique;size:16;comment:用户名"`
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
	userId := ctx.Value(global.ContextUserIDKey).(string)
	u.CreatedBy = userId
	u.UpdatedBy = userId

	return nil
}

func (u *SysUser) BeforeUpdate(tx *gorm.DB) (err error) {
	val := tx.Statement.Context.Value(global.ContextUserIDKey)
	userId, ok := val.(string)

	if ok && u.UpdatedBy != userId && u.UserId != "" {
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
	if u.UserId != "" {
		userId := tx.Statement.Context.Value(global.ContextUserIDKey).(string)
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
