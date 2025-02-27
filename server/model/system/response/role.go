package response

import (
	"github.com/imehc/do-exercise/server/model"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type RoleItem struct {
	model.IDWrapper      `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.RoleRequest  `json:",inline"`
	Dept                 DeptItem `json:"dept,omitzero"`
}

type RoleResponse struct {
	Data []RoleItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
