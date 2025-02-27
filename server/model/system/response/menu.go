package response

import (
	"github.com/imehc/do-exercise/server/model"
	commonRes "github.com/imehc/do-exercise/server/model/common/response"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
)

type MenuItem struct {
	model.IDWrapper      `json:",inline"`
	model.ControlWrapper `json:",inline"`
	request.MenuRequest  `json:",inline"`
	DataScope            string          `json:"data_scope"` // 数据权限
	IsSelect             bool            `json:"is_select"`  // 是否选中
	NoCache              bool            `json:"no_cache"`
	Params               string          `json:"params"`
	Paths                string          `json:"paths"`
	Apis                 []system.SysApi `json:"apis"`
	// Breadcrumb string `json:"breadcrumb"`	// 面包屑
}

type MenuListItem struct {
	MenuItem `json:",inline"`
	Children []MenuListItem `json:"children,inline"`
}

type MenuResponse struct {
	Data []MenuItem         `json:"data"`
	Meta commonRes.Paginate `json:"meta"`
}

type MenuTree struct {
	ID       int        `json:"id"`
	Label    string     `json:"label"`
	Children []MenuTree `json:"children"`
}
