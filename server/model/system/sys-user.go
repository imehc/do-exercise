package system

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
)

type SysUser struct {
	UserId   string `json:"id" gorm:"primarykey;column:id;type:varchar(32);comment:主键"`
	Username string `json:"username" gorm:"not null;size:16;uniqueIndex:idx_sys_user_username_tenant,priority:1,where:deleted_at IS NULL;comment:用户名"`
	Nickname string `json:"nickname" gorm:"comment:昵称"`
	// Email 邮箱唯一索引为条件索引（email <> ''），空邮箱不参与唯一性校验，
	// 从而同一租户下可存在多个空邮箱用户。
	Email    string    `json:"email" gorm:"uniqueIndex:idx_sys_user_email_tenant,priority:1,where:deleted_at IS NULL AND email <> '';comment:邮箱"`
	Avatar   string    `json:"avatar" gorm:"comment:头像"`
	Password string    `json:"-" gorm:"not null;comment:密码"`
	TenantId string    `json:"tenant_id" gorm:"column:tenant_id;type:varchar(32);not null;default:'';index;comment:租户ID;uniqueIndex:idx_sys_user_username_tenant,priority:2;uniqueIndex:idx_sys_user_email_tenant,priority:2"`
	Roles    []SysRole `json:"roles" gorm:"many2many:sys_user_role;comment:用户角色"`
	// MustChangePassword 标记账号是否需要在下次登录时强制修改密码。
	// 用于默认管理员口令的强制轮换：播种的默认口令是公开的，首次登录后必须改密。
	MustChangePassword bool `json:"must_change_password" gorm:"column:must_change_password;type:boolean;default:false;comment:是否需要在下次登录时强制修改密码"`
	// IsSuperAdmin 平台超级管理员标识（仅平台域 tenant_id=platform 的账号有效）。
	// 超管身份由标识决定而非账号名，改名/换账号只要标识为真即可，逻辑不依赖具体用户名。
	IsSuperAdmin bool `json:"is_super_admin" gorm:"column:is_super_admin;type:boolean;default:false;comment:是否平台超级管理员"`
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

	u.CreatedBy = model.CurrentUserID(tx)
	u.UpdatedBy = u.CreatedBy

	return nil
}

func (u *SysUser) BeforeUpdate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)

	if userId != "" && u.UpdatedBy != userId && u.UserId != "" {
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
		u.DeletedBy = model.CurrentUserID(tx)
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
