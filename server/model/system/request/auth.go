package request

import "github.com/imehc/do-exercise/server/utils"

type Login struct {
	Username  string `json:"username" binding:"required,min=4,max=8,alphanum,startWithLetter,containsLetter"`
	Password  string `json:"password" binding:"required,min=6,max=16,complexPassword"`
	Captcha   string `json:"captcha" binding:"required"`
	CaptchaId string `json:"captcha_id" binding:"required"`
} // @name LoginRequest

func (l Login) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"Username.required":        "用户名不能为空",
		"Username.min":             "用户名最少为4个字符",
		"Username.max":             "用户名最多为8个字符",
		"Username.alphanum":        "用户名只能包含字母和数字",
		"Username.startWithLetter": "用户名必须以字母开头",
		"Username.containsLetter":  "用户名必须包含至少一个字母",
		"Password.required":        "密码不能为空",
		"Password.min":             "密码最少为6个字符",
		"Password.max":             "密码最多为16个字符",
		"Password.complexPassword": "密码必须包含字母、数字和特殊字符",
		"Captcha.required":         "验证码不能为空",
		"CaptchaId.required":       "验证码ID不能为空",
	}
}

type RefreshToken struct {
	RefreshToken string `form:"refresh_token" binding:"required"`
} // @name RefreshTokenRequest

func (r RefreshToken) GetMessage() utils.ValidatorMessages {
	return utils.ValidatorMessages{
		"RefreshToken.required": "refreshToken不能为空",
	}
}
