package request

import (
	"github.com/imehc/do-exercise/server/model"
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

type PostRequest struct {
	Name                string `json:"name" binding:"required"` // 岗位名称
	Code                string `json:"code" binding:"required"` // 岗位编码
	DeptId              int    `json:"dept_id,omitzero"`        // 所属部门
	model.SortWrapper   `json:",inline"`
	model.StatusWrapper `json:",inline"`
	model.RemarkWrapper `json:",inline"`
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
	commonReq.QueryParams `json:",inline"`
	Name                  string `json:"name" form:"name"`
}
