package common

type Token struct {
	AccessToken       string `json:"access_token"`        // 访问令牌
	ExpireTime        int64  `json:"expire_time"`         // 访问令牌过期时间
	RefreshToken      string `json:"refresh_token"`       // 刷新令牌
	RefreshExpireTime int64  `json:"refresh_expire_time"` // 刷新令牌过期时间
}
