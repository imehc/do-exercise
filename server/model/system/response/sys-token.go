package response

import "time"

type SysTokenLogRsp struct {
	UserId              string    `json:"user_id"`
	Username            string    `json:"username"`
	AccessToken         string    `json:"access_token"`
	RefreshToken        string    `json:"refresh_token"`
	Disabled            bool      `json:"disabled"`
	AccessTokenCreated  time.Time `json:"access_token_created"`
	AccessTokenExpired  time.Time `json:"access_token_expired"`
	RefreshTokenCreated time.Time `json:"refresh_token_created"`
	RefreshTokenExpired time.Time `json:"refresh_token_expired"`
}
