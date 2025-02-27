package response

type Paginate struct {
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页显示数量
	Total    int64 `json:"total"`     // 总记录数
}
