package response

import "time"

type SysRoleResp struct {
	Id    uint               `json:"id"`
	Name  string             `json:"name"`
	Code  string             `json:"code"`
	Menus []SysMenuShortResp `json:"menus,omitzero"`
	// UserCount 当前持有该角色的未注销用户数。
	// 不加 omitzero：0 是有意义的取值（「还没有人用」），省略会让前端分不清
	// 「没人用」和「后端没给」。
	UserCount int64     `json:"user_count"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	CreatedBy string    `json:"created_by,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
	UpdatedBy string    `json:"updated_by,omitzero"`
}

type SysRoleShortResp struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
