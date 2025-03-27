package common

type Login struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required,min=6,max=16,complexPassword"`
}

type LoginReq struct {
	Login
	// TODO: 验证码相关
}
