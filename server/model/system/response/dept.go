package response

import (
	"github.com/imehc/do-exercise/server/global"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type DeptItem struct {
	global.Model
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
