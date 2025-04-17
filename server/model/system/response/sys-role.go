package response

import "time"

type SysRoleResp struct {
	Id        uint               `json:"id"`
	Name      string             `json:"name"`
	Code      string             `json:"code"`
	Menus     []SysMenuShortResp `json:"menus"`
	CreatedAt time.Time          `json:"created_at"`
	CreatedBy int64              `json:"created_by"`
	UpdatedAt time.Time          `json:"updated_at"`
	UpdatedBy int64              `json:"updated_by"`
}
