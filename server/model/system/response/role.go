package response

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type RoleItem struct {
	common.IDWrapper          `json:",inline"`
	model.ControlWrapper      `json:",inline"`
	request.CreateRoleRequest `json:",inline"`
	Menus                     []MenuItem `json:"menus"`
	Depts                     []DeptItem `json:"depts,omitzero"`
	Apis                      []ApiItem  `json:"apis"`
	IsAdmin                   bool       `json:"is_admin"`
	DataScope                 uint       `json:"data_scope"`
}

type RoleResponse struct {
	Data []RoleItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
