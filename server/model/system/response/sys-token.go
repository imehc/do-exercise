package response

import "time"

type SysTokenLogRsp struct {
	UserId       int64     `json:"user_id"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Disabled     bool      `json:"disabled"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}
