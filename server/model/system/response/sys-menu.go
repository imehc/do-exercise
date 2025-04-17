package response

import "time"

type SysMenuShortResp struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type SysMenuResp struct {
	Id         uint      `json:"id"`
	Name       string    `json:"name"`
	ParentId   *uint     `json:"parent_id"`
	Permission *string   `json:"permission,omitzero"`
	Icon       string    `json:"icon,omitzero"`
	Type       uint8     `json:"type"`
	Route      string    `json:"route,omitzero"`
	Component  string    `json:"component,omitzero"`
	Sort       uint      `json:"sort"`
	Visible    bool      `json:"visible"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  int64     `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
	UpdatedBy  int64     `json:"updated_by"`
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
