package response

import (
	"github.com/imehc/do-exercise/server/model"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type PostItem struct {
	model.IDWrapper
	model.ControlWrapper
	request.PostRequest
	Dept *DeptItem `json:"dept,omitempty"`
}

type PostResponse struct {
	Data []PostItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
