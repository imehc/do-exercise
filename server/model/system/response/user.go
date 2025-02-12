package response

import (
	"github.com/imehc/do-exercise/server/global"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type UserItem struct {
	global.Model
	request.UserItem
	Dept    DeptItem `json:"dept"`
	DeptIds []int    `json:"dept_ids"`
	PostIds []int    `json:"post_ids"`
	RoleIds []int    `json:"role_ids"`
}

type UserResponse struct {
	Data []UserItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
