package common

// MaxPageSize 单页最大记录数，防止 page_size 过大导致内存与数据库压力
const MaxPageSize = 100

// Pagination 分页请求参数
type Pagination struct {
	Page     int `json:"page" form:"page" binding:"required,min=1"`                   // 当前页码
	PageSize int `json:"page_size" form:"page_size" binding:"required,min=1,max=100"` // 每页数量
}

// Normalize 将分页参数规整到合法范围，避免越界请求产生非法查询
func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}
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
