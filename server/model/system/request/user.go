package request

import (
	"github.com/imehc/do-exercise/server/model"
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

type UserItem struct {
	Avatar   string `json:"avatar"`                                                                          // 头像路径
	Nickname string `json:"nickname"`                                                                        // 昵称
	DeptId   int    `json:"dept_id" binding:"required"`                                                      // 归属部门
	Phone    string `json:"phone"`                                                                           // 手机号
	Email    string `json:"email"`                                                                           // 邮箱
	Username string `json:"username" binding:"required,min=4,max=8,alphanum,startWithLetter,containsLetter"` // 用户名
	Sex      int    `json:"sex"`                                                                             // 性别
	PostId   int    `json:"post_id" binding:"required"`                                                      // 岗位
	RoleId   int    `json:"role_id" binding:"required"`                                                      // 角色

	model.RemarkWrapper
	model.StatusWrapper
}

type UserRequest struct {
	UserItem
	Password string `json:"password" binding:"required,min=6,max=16,complexPassword"` // 密码
}

func (u *UserRequest) EncryptPassword() {
	u.Password = utils.BcryptHash(u.Password)
}

func (u UserRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Username.required":        "用户名不能为空",
		"Username.min":             "用户名最少为4个字符",
		"Username.max":             "用户名最多为8个字符",
		"Username.alphanum":        "用户名只能包含字母和数字",
		"Username.startWithLetter": "用户名必须以字母开头",
		"Username.containsLetter":  "用户名必须包含至少一个字母",
		"Password.required":        "密码不能为空",
		"Password.min":             "密码最少为6个字符",
		"Password.max":             "密码最多为16个字符",
		"Password.complexPassword": "密码必须包含字母、数字和特殊字符",
		"DeptId.required":          "归属部门不能为空",
		"PostId.required":          "岗位不能为空",
		"RoleId.required":          "角色不能为空",
	}
}

type UserParam struct {
	UserId int `json:"user_id" binding:"required"`
}

type UserQueryParams struct {
	commonReq.QueryParams
	Name string `json:"name" form:"name"`
}
