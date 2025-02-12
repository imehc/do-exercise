package request

import (
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

type DeptRequest struct {
	ParentId *int   `json:"parent_id" binding:"required"` // 上级部门
	Name     string `json:"name" binding:"required"`      // 部门名称
	Sort     int    `json:"sort"`                         // 排序
	Leader   string `json:"leader" binding:"required"`    // 负责人
	Phone    string `json:"phone"`                        // 手机号
	Email    string `json:"email"`                        // 邮箱
	Status   int    `json:"status"`                       // 状态
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
	commonReq.QueryParams
	Name string `json:"name" form:"name"`
}
