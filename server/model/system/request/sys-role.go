package request

import "github.com/imehc/do-exercise/server/model/common"

type CreateSysRoleReq struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required,alphanum"`
	MenuIds []uint `json:"menu_ids"`
}

type UpdateSysRoleReq struct {
	Name    string `json:"name" binding:"required"`
	MenuIds []uint `json:"menu_ids"`
}

type QuerySysRoleReq struct {
	common.Pagination
	Name string `json:"name" form:"name"` // 角色名称
	Code string `json:"code" form:"code"` // 角色编码
}
