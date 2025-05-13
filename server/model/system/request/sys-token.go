package request

type SysTokenDeleteReq struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type SysTokenModityStatusReq struct {
	AccessToken string `json:"access_token" binding:"required"`
	Disabled    bool   `json:"disabled" binding:"required"`
}
