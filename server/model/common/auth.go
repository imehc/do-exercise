package common

type Captcha struct {
	CaptchaId string `json:"captcha_id"` // 验证码ID
	PicPath   string `json:"pic_path"`   // 验证码图片路径
}

type Login struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required"`
}

type LoginReq struct {
	Login     `json:",inline"`
	Captcha   string `json:"captcha" binding:"required,min=1,max=8"`
	CaptchaId string `json:"captcha_id" binding:"required"`
	PublicKey string `json:"public_key"`
}
