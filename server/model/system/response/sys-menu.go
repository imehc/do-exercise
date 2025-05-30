package response

import "time"

type SysMenuShortResp struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type SysMenuResp struct {
	Id         uint         `json:"id"`
	Name       string       `json:"name"`
	ParentId   *uint        `json:"parent_id"`
	Permission *string      `json:"permission,omitzero"`
	Icon       string       `json:"icon,omitzero"`
	Type       uint8        `json:"type"`
	Route      string       `json:"route,omitzero"`
	Component  string       `json:"component,omitzero"`
	Sort       uint         `json:"sort"`
	Visible    bool         `json:"visible"`
	Apis       []SysApiResp `json:"apis,omitzero"`
	CreatedAt  time.Time    `json:"created_at,omitzero"`
	CreatedBy  string       `json:"created_by,omitzero"`
	UpdatedAt  time.Time    `json:"updated_at,omitzero"`
	UpdatedBy  string       `json:"updated_by,omitzero"`
}

type SysMenuTreeResp struct {
	Id         uint              `json:"id"`
	Name       string            `json:"name"`
	ParentId   *uint             `json:"parent_id"`
	Permission *string           `json:"permission,omitzero"`
	Icon       string            `json:"icon,omitzero"`
	Type       uint8             `json:"type"`
	Route      string            `json:"route,omitzero"`
	Component  string            `json:"component,omitzero"`
	Sort       uint              `json:"sort"`
	Visible    bool              `json:"visible"`
	Children   []SysMenuTreeResp `json:"children"`
}
