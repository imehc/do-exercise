package response

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type PostItem struct {
	common.IDWrapper     `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.PostRequest  `json:",inline"`
	Dept                 DeptItem `json:"dept,omitzero"`
}

type PostResponse struct {
	Data []PostItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}
