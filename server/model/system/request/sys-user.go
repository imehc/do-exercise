package request

type CreateSysUserReq struct {
	Username string `json:"username" validate:"required,alphanum,min=2,max=10"`
	Password string `json:"password" validate:"required,alphanum,min=6,max=16"`
	Phone    string `json:"phone"`
	Email    string `json:"email" validate:"email"`
	Avatar   string `json:"avatar"`
}
