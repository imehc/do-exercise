package response

import "time"

type SysRoleResp struct {
	Id        uint               `json:"id"`
	Name      string             `json:"name"`
	Code      string             `json:"code"`
	Menus     []SysMenuShortResp `json:"menus,omitzero"`
	CreatedAt time.Time          `json:"created_at,omitzero"`
	CreatedBy string             `json:"created_by,omitzero"`
	UpdatedAt time.Time          `json:"updated_at,omitzero"`
	UpdatedBy string             `json:"updated_by,omitzero"`
}

type SysRoleShortResp struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
