package request

import (
	commonReq "github.com/imehc/do-exercise/server/model/common/request"
	"github.com/imehc/do-exercise/server/utils"
)

type ApiRequest struct {
	Handle string `json:"handle" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Path   string `json:"path" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Action string `json:"action" binding:"required"`
}

func (p ApiRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Handle.required": "handle不能为空",
		"Title.required":  "标题不能为空",
		"Path.required":   "地址不能为空",
		"Type.required":   "接口类型不能为空",
		"Action.required": "请求类型不能为空",
	}
}

type ApiParam struct {
	ApiId int `json:"api_id" binding:"required"`
}

type ApiQueryParams struct {
	commonReq.QueryParams `json:",inline"`
	Name                  string `json:"name" form:"name"`
}
