package response

type Role struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type UserProfile struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Roles    []Role `json:"roles"`
}

type UserMenu struct {
	Id         uint    `json:"id"`
	Name       string  `json:"name"`
	ParentId   *uint   `json:"parent_id"`
	Permission *string `json:"permission,omitzero"`
	Icon       string  `json:"icon,omitzero"`
	Type       uint8   `json:"type"`
	Route      string  `json:"route,omitzero"`
	Component  string  `json:"component,omitzero"`
}
