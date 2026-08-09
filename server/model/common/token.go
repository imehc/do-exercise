package common

type Token struct {
	AccessToken       string `json:"access_token"`        // 访问令牌
	ExpireTime        int64  `json:"expire_time"`         // 访问令牌过期时间
	RefreshToken      string `json:"refresh_token"`       // 刷新令牌
	RefreshExpireTime int64  `json:"refresh_expire_time"` // 刷新令牌过期时间
	// MustChangePassword 标记该账号仍需强制修改密码，前端据此强制跳转改密页
	MustChangePassword bool `json:"must_change_password"`
}
