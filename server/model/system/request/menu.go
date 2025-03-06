package request

import (
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/utils"
)

type MenuRequest struct {
	ParentId   *uint  `json:"parent_id" binding:"required"` // 上级菜单ID
	Name       string `json:"name"`                         // 路由名称
	Icon       string `json:"icon"`                         // 菜单图标
	Type       string `json:"type" binding:"required"`      // 菜单类型（M目录 C菜单 F按钮）
	Action     string `json:"action"`                       // 路由地址
	Visible    bool   `json:"visible"`                      // 显示状态
	Title      string `json:"title" binding:"required"`     // 菜单标题
	Route      string `json:"route"`                        // 路由地址
	Permission string `json:"permission"`                   // 权限标识
	ApiIds     []uint `json:"api_ids,omitzero"`

	model.SortWrapper `json:",inline"`
}

func (m MenuRequest) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"ParentId.required": "上级菜单不能为空",
		"Title.required":    "菜单标题不能为空",
		"Type.required":     "菜单类型不能为空",
	}
}

type MenuParam struct {
	MenuId int `json:"menu_id" binding:"required"`
}

type MenuQueryParams struct {
	Name string `json:"name" form:"name"`
}
