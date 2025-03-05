package response

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type DeptItem struct {
	common.IDWrapper     `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.DeptRequest  `json:",inline"`
}

type DeptResponse struct {
	DeptItem `json:",inline"`
	Children []DeptResponse `json:"children"`
}

type DeptTree struct {
	ID       int        `json:"id"`
	Label    string     `json:"label"`
	Children []DeptTree `json:"children"`
}
