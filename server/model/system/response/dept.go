package response

import (
	"github.com/imehc/do-exercise/server/model"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type DeptItem struct {
	model.IDWrapper
	model.ControlWrapper
	request.DeptRequest
}

type DeptResponse struct {
	Data []DeptItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}

type DeptTree struct {
	ID       int        `json:"id"`
	Label    string     `json:"label"`
	Children []DeptTree `json:"children"`
}
