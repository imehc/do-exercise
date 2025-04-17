package request

type CreateSysUserReq struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required,min=6,max=16,complexPassword"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"email"`
	Avatar   string `json:"avatar"`
}

type UpdateSysUserReq struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
