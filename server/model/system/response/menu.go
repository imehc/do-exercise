package response

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type MenuItem struct {
	common.IDWrapper     `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.MenuRequest  `json:",inline"`
	DataScope            string       `json:"data_scope"` // 数据权限
	IsSelect             bool         `json:"is_select"`  // 是否选中
	NoCache              bool         `json:"no_cache"`
	Params               string       `json:"params"`
	Route                string       `json:"route"`
	Icon                 string       `json:"icon"`
	Apis                 []system.Api `json:"apis,omitzero"`
	Children             []MenuItem   `json:"children,omitzero"`
	// Breadcrumb string `json:"breadcrumb"`	// 面包屑
}

type MenuListItem struct {
	MenuItem `json:",inline"`
	Children []MenuListItem `json:"children,inline"`
}

type MenuCompact struct {
	ID       int           `json:"id"`
	Label    string        `json:"label"`
	Route    string        `json:"route"`
	Icon     string        `json:"icon"`
	Children []MenuCompact `json:"children"`
}

type MenuTree struct {
	ID       int        `json:"id"`
	Label    string     `json:"label"`
	Children []MenuTree `json:"children"`
}
