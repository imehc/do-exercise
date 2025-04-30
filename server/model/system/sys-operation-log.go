package system

import (
	"time"

	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

type SysOperationLog struct {
	model.IdWrapper

	Username       string         `json:"username" gorm:"size:16;comment:请求用户"`
	Ip             string         `json:"ip" gorm:"size:32;comment:IP地址"`
	IsInternalIP   bool           `json:"is_internal_ip" gorm:"default:false;comment:是否是内网IP"`
	Address        string         `json:"address" gorm:"size:128;comment:地址"`
	UserAgent      string         `json:"user_agent" gorm:"size:255;comment:ua"`
	Borwser        string         `json:"browser" gorm:"size:128;comment:浏览器"`
	BorwserVersion string         `json:"browser_version" gorm:"size:128;comment:浏览器版本"`
	IsMobile       bool           `json:"is_mobile" gorm:"default:false;comment:是否是手机"`
	IsBot          bool           `json:"is_bot" gorm:"default:false;comment:是否是机器人"`
	Os             string         `json:"os" gorm:"size:128;comment:操作系统"`
	Path           string         `json:"path" gorm:"size:128;comment:请求路径"`
	Method         string         `json:"method" gorm:"size:16;comment:请求方法"`
	Success        bool           `json:"success" gorm:"default:false;comment:是否成功"`
	Code           int            `json:"code" gorm:"comment:错误码"`
	Message        string         `json:"message" gorm:"size:255;comment:错误信息"`
	Params         string         `json:"params" gorm:"type:text;comment:请求参数"`
	Body           string         `json:"body" gorm:"type:text;comment:请求体"`
	Result         string         `json:"result" gorm:"type:text;comment:响应结果"`
	StartTime      time.Time      `json:"start_time" gorm:"comment:请求开始时间"`
	EndTime        time.Time      `json:"end_time" gorm:"comment:请求结束时间"`
	Latency        int64          `json:"latency" gorm:"comment:耗时"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt      time.Time      `json:"updated_at,omitzero" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index;comment:删除时间"`
}
