package response

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type MenuItem struct {
	model.IDWrapper      `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.MenuRequest  `json:",inline"`
	DataScope            string       `json:"data_scope"` // 数据权限
	IsSelect             bool         `json:"is_select"`  // 是否选中
	NoCache              bool         `json:"no_cache"`
	Params               string       `json:"params"`
	Paths                string       `json:"paths"`
	Apis                 []system.Api `json:"apis"`
	Children             []MenuItem   `json:"children"`
	// Breadcrumb string `json:"breadcrumb"`	// 面包屑
}

type MenuListItem struct {
	MenuItem `json:",inline"`
	Children []MenuListItem `json:"children,inline"`
}

type MenuTree struct {
	ID       int        `json:"id"`
	Label    string     `json:"label"`
	Children []MenuTree `json:"children"`
}
