package request

import (
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

type RoleRequest struct {
}

type CreateRoleRequest struct {
	RoleRequest `json:",inline"`
}

func (d CreateRoleRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type UpdateRoleRequest struct {
	RoleRequest `json:",inline"`
}

func (d UpdateRoleRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type UpdateRoleDataScope struct{}

func (d UpdateRoleDataScope) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type UpdateMenuDataScope struct{}

func (d UpdateMenuDataScope) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type UpdateApiDataScope struct{}

func (d UpdateApiDataScope) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type RoleParam struct {
	RoleId int `json:"role_id" binding:"required"`
}

type RoleQueryParams struct {
	commonReq.QueryParams `json:",inline"`
	Name                  string `json:"name" form:"name"`
}
