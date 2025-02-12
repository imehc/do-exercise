package request

import "github.com/imehc/do-exercise/server/utils"

type QueryParams struct {
	Page     int `json:"page" form:"page" binding:"min=1"`                   // 分页页码
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"` // 分页每页大小
}

func (q QueryParams) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Page.min":     "分页页码最小为1",
		"PageSize.min": "分页每页大小最小为1",
		"PageSize.max": "分页每页大小最大为100",
	}
}
