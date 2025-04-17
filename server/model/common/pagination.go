package common

// Pagination 分页请求参数
type Pagination struct {
	Page     int `json:"page" form:"page" binding:"required,min=1"`           // 当前页码
	PageSize int `json:"page_size" form:"page_size" binding:"required,min=1"` // 每页数量
}

// PageMeta 分页元信息
type PageMeta struct {
	Total    int64 `json:"total"`     // 总记录数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页数量
}

// PageResult 通用分页响应
type PageResult[T any] struct {
	Data []T      `json:"data"` // 数据列表
	Meta PageMeta `json:"meta"` // 分页元信息
}
