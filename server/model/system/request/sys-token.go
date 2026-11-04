package request

type SysTokenDeleteReq struct {
	SessionId string `json:"session_id" binding:"required"`
}

type SysTokenModityStatusReq struct {
	SessionId string `json:"session_id" binding:"required"`
	Disabled  *bool  `json:"disabled" binding:"required"`
}
