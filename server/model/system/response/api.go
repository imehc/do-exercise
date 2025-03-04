package response

import (
	"github.com/imehc/do-exercise/server/model"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type ApiItem struct {
	model.IDWrapper      `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.ApiRequest   `json:",inline"`
}

type ApiResponse struct {
	Data []ApiItem          `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
