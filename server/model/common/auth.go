package common

type Captcha struct {
	CaptchaId string `json:"captcha_id"` // 验证码ID
	PicPath   string `json:"pic_path"`   // 验证码图片路径
}

type Login struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required"`
	// TenantId 多租户模式下显式指定登录租户；为空时按用户名自动解析：
	// 单租户模式用默认租户，多租户模式按用户名归属的启用租户唯一/多选处理。
	TenantId string `json:"tenant_id"`
}

type LoginReq struct {
	Login     `json:",inline"`
	Captcha   string `json:"captcha" binding:"required,min=1,max=8"`
	CaptchaId string `json:"captcha_id" binding:"required"`
	PublicKey string `json:"public_key"`
}

// LoginSession 多租户登录会话：用户名与密码验证通过，但账号归属多个启用租户时，
// 后端签发一次性会话，等待前端选择要进入的租户后再发正式 token。
type LoginSession struct {
	Username string   `json:"username"`
	Tenants  []string `json:"tenants"`
}

// SelectTenantReq 多租户登录选择租户请求
type SelectTenantReq struct {
	LoginSessionId string `json:"login_session_id" binding:"required"`
	TenantId       string `json:"tenant_id" binding:"required"`
}

// SwitchTenantReq 登录后切换租户请求
type SwitchTenantReq struct {
	TenantId string `json:"tenant_id" binding:"required"`
}

type ResetPasswordReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type LoginWithEmailReq struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
