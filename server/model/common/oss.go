package common

type OssReq struct {
	FileName string `json:"file_name" form:"file_name" binding:"required"` // 文件名
}

type OssRes struct {
	PutObjectUrl string `json:"put_object_url"` // 上传地址
	GetObjectUrl string `json:"get_object_url"` // 下载地址
	Expires      int    `json:"expires"`        // 过期时间戳
}
