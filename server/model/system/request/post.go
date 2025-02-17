package request

import (
	"github.com/imehc/do-exercise/server/model"
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

type PostRequest struct {
	Name   string `json:"name" binding:"required"` // 岗位名称
	Code   string `json:"code" binding:"required"` // 岗位编码
	DeptId *int   `json:"dept_id,omitempty"`       // 所属部门
	// DeptId   int    `json:"dept_id" binding:"required"`   // 所属部门

	model.SortWrapper
	model.StatusWrapper
	model.RemarkWrapper
}

func (p PostRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Name.required":     "岗位名称不能为空",
		"PostCode.required": "岗位编码不能为空",
		// "DeptId.required":   "所属部门不能为空",
	}
}

type PostParam struct {
	PostId int `json:"post_id" binding:"required"`
}

type PostQueryParams struct {
	commonReq.QueryParams
	Name string `json:"name" form:"name"`
}
