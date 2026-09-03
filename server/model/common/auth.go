package common

type Captcha struct {
	CaptchaId string `json:"captcha_id"` // 验证码ID
	PicPath   string `json:"pic_path"`   // 验证码图片路径
}

type Login struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=10,startWithLetter,containsLetter"`
	Password string `json:"password" binding:"required"`
	// TenantId 显式指定登录租户；为空时按用户名归属的启用租户自动解析：
	// 唯一归属直接登录，多个归属返回候选列表要求选择。
	TenantId string `json:"tenant_id"`
	// TenantCode 登录页可填的租户编码，等价于显式指定租户（TenantId 优先）。
	// 填了就只在该租户下验证一次口令，避免同名账号跨租户逐行比对的成本放大；
	// 代价是本次会话被钉在该租户上，需要跨租户切换请留空。
	// 校验口径与 CreateSysTenantReq.Code 一致（alphanum,min=2,max=32）：编码只能在创建时设置，
	// 所以库里不可能存在不满足该规则的编码，提前拒绝可以省掉一次注定查不到的查询。
	TenantCode string `json:"tenant_code" binding:"omitempty,alphanum,min=2,max=32"`
}

type LoginReq struct {
	Login     `json:",inline"`
	Captcha   string `json:"captcha" binding:"required,min=1,max=8"`
	CaptchaId string `json:"captcha_id" binding:"required"`
	PublicKey string `json:"public_key"`
}

// LoginSession 待选择租户的一次性登录会话。
//
// 两条来源，收敛方式不同：
//   - 口令登录：口令已在 sys_user 逐行验证过，用 Username 即可在目标租户下重新定位账号；
//   - 邮箱登录：同一邮箱在不同租户下可以是不同的用户名，Username 不足以定位，
//     因此记录本次验证码绑定的候选账号 ID（UserIds），选租户后按 ID 收敛。
//
// 两者都只能进入 Tenants 里列出的租户，客户端无法扩大授权范围。
type LoginSession struct {
	Username string   `json:"username"`
	Tenants  []string `json:"tenants"`
	UserIds  []string `json:"user_ids"`
}

// SelectTenantReq 多租户登录选择租户请求
type SelectTenantReq struct {
	LoginSessionId string `json:"login_session_id" binding:"required"`
	TenantId       string `json:"tenant_id" binding:"required"`
}

// SwitchTenantReq 登录后切换租户请求
type SwitchTenantReq struct {
	TenantId  string `json:"tenant_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"`
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
