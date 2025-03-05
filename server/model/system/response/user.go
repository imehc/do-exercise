package response

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type UserItem struct {
	common.IDWrapper     `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.UserItem     `json:",inline"`
	Dept                 DeptItem `json:"dept,omitzero"`
	Post                 PostItem `json:"post,omitzero"`
}

type UserResponse struct {
	Data []UserItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
