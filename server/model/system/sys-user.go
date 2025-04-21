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
	return nil
}
