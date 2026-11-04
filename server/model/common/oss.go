package common

type OssReq struct {
	FileName string `json:"file_name" form:"file_name" binding:"required,max=255"`
	FileSize int64  `json:"file_size" form:"file_size" binding:"required,gt=0,lte=104857600"`
	Purpose  string `json:"purpose" form:"purpose" binding:"omitempty,oneof=avatar private"`
}
type OssRes struct {
	UploadUrl    string            `json:"upload_url"`
	FormFields   map[string]string `json:"form_fields"`
	ObjectKey    string            `json:"object_key"`
	GetObjectUrl string            `json:"get_object_url,omitempty"`
	Expires      int               `json:"expires"`
}
type OssDownloadReq struct {
	ObjectKey string `form:"object_key" binding:"required,max=512"`
}
type OssDownloadRes struct {
	Url     string `json:"url"`
	Expires int    `json:"expires"`
}

type OssObjectReq struct {
	ObjectKey string `json:"object_key" binding:"required,max=512"`
}
type OssFileRes struct {
	ObjectKey    string `json:"object_key"`
	Size         int64  `json:"size"`
	State        string `json:"state"`
	GetObjectUrl string `json:"get_object_url,omitempty"`
}
type OssUsageRes struct {
	ChargedBytes    int64 `json:"charged_bytes"`
	MaxStorageBytes int64 `json:"max_storage_bytes"`
}
