package response

// ValidationError 表示验证错误的响应结构
type ValidationError struct {
	Type    string             `json:"type"`
	Message string             `json:"message"`
	Details []ValidationDetail `json:"details,omitzero"`
}

// ValidationDetail 表示具体的验证错误信息
type ValidationDetail struct {
	Field    string   `json:"field"`
	Messages []string `json:"messages"`
}
