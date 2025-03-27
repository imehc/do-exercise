package request

type CreateSysUserReq struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required,min=6,max=16,complexPassword"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"email"`
	Avatar   string `json:"avatar"`
}

type LoginReq struct {
	Login
	// TODO: 验证码相关
}

type Login struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required,min=6,max=16,complexPassword"`
}
