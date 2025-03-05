package system

import (
	"github.com/imehc/do-exercise/server/model"
)

type User struct {
	UserId   uint   `gorm:"primarykey;autoIncrement;comment:主键ID" json:"user_id"`
	Username string `gorm:"uniqueIndex;size:16;comment:用户名" json:"username"`
	Password string `gorm:"not null;size:128;comment:密码" json:"password"`
	Email    string `gorm:"size:128;comment:邮箱" json:"email"`
	Nickname string `gorm:"size:128;comment:昵称" json:"nickname"`
	Phone    string `gorm:"size:11;comment:联系电话" json:"phone"`
	Avatar   string `gorm:"size:255;comment:头像" json:"avatar"`
	Sex      int    `gorm:"size:255;comment:性别" json:"sex"`
	RoleId   uint   `gorm:"size:20;comment:角色ID" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleId" json:"role"` // 对应的详情数据
	DeptId   uint   `gorm:"size:20;comment:部门" json:"dept_id"`
	Dept     Dept   `gorm:"foreignKey:DeptId" json:"dept"` // 对应的详情数据
	PostId   uint   `gorm:"size:20;comment:岗位" json:"post_id"`
	Post     Post   `gorm:"foreignKey:PostId" json:"post"` // 对应的详情数据

	model.RemarkWrapper
	model.StatusWrapper
	model.ControlWrapper
}

func (*User) TableName() string {
	return "sys_user"
}
