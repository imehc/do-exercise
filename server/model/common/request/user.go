package request

type Email struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}

type EmailCache struct {
	Code   string `json:"code"`
	UserId int64  `json:"user_id"`
}

type BindEmailReq struct {
	Id    int64  `json:"id"`
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type UserResetPasswordReq struct {
	Id        int64  `json:"id"`
	Email     string `json:"email" binding:"required,email"`
	Code      string `json:"code" binding:"required"`
	Password  string `json:"password" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"`
}

type UserModifyPasswordReq struct {
	Id          int64  `json:"id"`
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PublicKey   string `json:"public_key" binding:"required"`
}

type UserModifyProfileReq struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
