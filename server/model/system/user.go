package system

import (
	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

type User struct {
	global.Model
	global.ControlBy

	Username string `gorm:"uniqueIndex;size:16;comment:用户名" json:"username"`
	Password string `gorm:"not null;size:128;comment:密码" json:"password"`
	Email    string `json:"email" gorm:"size:128;comment:邮箱"`
	Nickname string `json:"nickname" gorm:"size:128;comment:昵称"`
	Phone    string `json:"phone" gorm:"size:11;comment:联系电话"`
	Avatar   string `json:"avatar" gorm:"size:255;comment:头像"`
	Sex      int    `json:"sex" gorm:"size:255;comment:性别"`
	Remark   string `json:"remark" gorm:"size:255;comment:备注"`
	Status   int    `json:"status" gorm:"size:4;comment:状态"`
	RoleId   int    `json:"role_id" gorm:"size:20;comment:角色ID"`
	DeptId   int    `json:"dept_id" gorm:"size:20;comment:部门"`
	Dept     Dept   `json:"dept" gorm:"foreignKey:DeptId"` // 对应的详情数据
	PostId   int    `json:"post_id" gorm:"size:20;comment:岗位"`
	DeptIds  []int  `json:"dept_ids" gorm:"-"`
	PostIds  []int  `json:"post_ids" gorm:"-"`
	RoleIds  []int  `json:"role_ids" gorm:"-"`
	// TODO: 更多信息使用单独一个表
}

func (u *User) AfterFind(_ *gorm.DB) error {
	u.DeptIds = []int{u.DeptId}
	u.PostIds = []int{u.PostId}
	u.RoleIds = []int{u.RoleId}
	return nil
}

func (*User) TableName() string {
	return "sys_user"
}
