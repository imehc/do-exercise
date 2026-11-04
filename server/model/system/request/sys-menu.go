package request

type CreateSysMenuReq struct {
	Name       string  `json:"name" binding:"required"`
	I18nKey    *string `json:"i18n_key" binding:"omitempty,max=128"`
	ParentId   *uint   `json:"parent_id" binding:"required"`
	Permission *string `json:"permission" binding:"omitempty,uniqueName"`
	Icon       string  `json:"icon"`
	Type       uint8   `json:"type" binding:"required,menuType"`
	Route      string  `json:"route"`
	Component  string  `json:"component"`
	Sort       uint    `json:"sort"`
	Visible    bool    `json:"visible"`
	Scope      string  `json:"scope" binding:"omitempty,oneof=platform tenant both"`
	IsSystem   bool    `json:"is_system"`
	ApiIds     []uint  `json:"api_ids"`
}

type UpdateSysMenuReq struct {
	ApiIds           []uint `json:"api_ids"`
	CreateSysMenuReq `json:",inline"`
}
