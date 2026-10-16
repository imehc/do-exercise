package response

import "time"

type SysMenuShortResp struct {
	Id      uint    `json:"id"`
	Name    string  `json:"name"`
	I18nKey *string `json:"i18n_key,omitzero"`
}

type SysMenuResp struct {
	Id         uint         `json:"id"`
	Name       string       `json:"name"`
	I18nKey    *string      `json:"i18n_key,omitzero"`
	ParentId   *uint        `json:"parent_id"`
	Permission *string      `json:"permission,omitzero"`
	Icon       string       `json:"icon,omitzero"`
	Type       uint8        `json:"type"`
	Route      string       `json:"route,omitzero"`
	Component  string       `json:"component,omitzero"`
	Sort       uint         `json:"sort"`
	Visible    bool         `json:"visible"`
	Scope      string       `json:"scope"`
	IsSystem   bool         `json:"is_system"`
	Apis       []SysApiResp `json:"apis,omitzero"`
	CreatedAt  time.Time    `json:"created_at,omitzero"`
	CreatedBy  string       `json:"created_by,omitzero"`
	UpdatedAt  time.Time    `json:"updated_at,omitzero"`
	UpdatedBy  string       `json:"updated_by,omitzero"`
}

// SysMenuApiBrief 菜单绑定接口的精简视图。
// 角色授权页需要在勾选一条权限时就地说明它实际放开哪些接口，
// 但树接口一次返回全量菜单，挂完整的 SysApiResp（含时间戳）会明显放大响应体。
type SysMenuApiBrief struct {
	Id          uint   `json:"id"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description,omitzero"`
}

type SysMenuTreeResp struct {
	Id         uint              `json:"id"`
	Name       string            `json:"name"`
	I18nKey    *string           `json:"i18n_key,omitzero"`
	ParentId   *uint             `json:"parent_id"`
	Permission *string           `json:"permission,omitzero"`
	Icon       string            `json:"icon,omitzero"`
	Type       uint8             `json:"type"`
	Route      string            `json:"route,omitzero"`
	Component  string            `json:"component,omitzero"`
	Sort       uint              `json:"sort"`
	Visible    bool              `json:"visible"`
	Scope      string            `json:"scope"`
	IsSystem   bool              `json:"is_system"`
	Apis       []SysMenuApiBrief `json:"apis,omitzero"`
	Children   []SysMenuTreeResp `json:"children"`
}
