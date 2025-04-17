package request

type CreateSysMenuReq struct {
	Name       string  `json:"name" binding:"required"`
	ParentId   *uint   `json:"parent_id" binding:"required"`
	Permission *string `json:"permission" binding:"omitempty,uniqueName"`
	Icon       string  `json:"icon"`
	Type       uint8   `json:"type" binding:"required,menuType"`
	Route      string  `json:"route"`
	Component  string  `json:"component"`
	Sort       uint    `json:"sort"`
	Visible    bool    `json:"visible"`
}

type UpdateSysMenuReq struct {
	Id               uint `json:"id"`
	CreateSysMenuReq `json:",inline"`
}
