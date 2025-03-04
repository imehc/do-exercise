package request

import (
	"github.com/imehc/do-exercise/server/model"
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

type RoleRequest struct {
	model.SortWrapper
	model.RemarkWrapper
	model.StatusWrapper
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Key         string `json:"key" binding:"required"`
	RoleRequest `json:",inline"`
}

func (d CreateRoleRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Name.required": "角色名称不能为空",
		"Key.required":  "角色标识不能为空",
	}
}

type UpdateRoleRequest struct {
	RoleRequest `json:",inline"`
}

func (d UpdateRoleRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type UpdateRoleDataScope struct {
	DataScope uint  `json:"data_scope" binding:"required"`
	DeptIds   []int `json:"dept_ids"`
}

func (d UpdateRoleDataScope) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"DataScope.required": "数据范围不能为空",
	}
}

type UpdateMenuDataScope struct {
	MenuIds []int `json:"menu_ids"`
}

func (d UpdateMenuDataScope) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{}
}

type UpdateApiDataScope struct {
	ApiIds []int `json:"api_ids"`
}

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
