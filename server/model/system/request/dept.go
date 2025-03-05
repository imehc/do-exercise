package request

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/utils"
)

type DeptRequest struct {
	ParentId            *uint  `json:"parent_id,omitzero" binding:"required"` // 上级部门
	Name                string `json:"name" binding:"required"`               // 部门名称
	Leader              string `json:"leader" binding:"required"`             // 负责人
	Phone               string `json:"phone"`                                 // 手机号
	Email               string `json:"email"`                                 // 邮箱
	model.SortWrapper   `json:",inline"`
	model.StatusWrapper `json:",inline"`
}

func (d DeptRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"ParentId.required": "上级部门不能为空",
		"Name.required":     "部门名称不能为空",
		"Leader.required":   "负责人不能为空",
	}
}

type DeptParam struct {
	DeptId int `json:"dept_id" binding:"required"`
}

type DeptQueryParams struct {
	Name string `json:"name" form:"name"`
}
