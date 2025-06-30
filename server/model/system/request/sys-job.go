package request

import "github.com/imehc/do-exercise/server/model/common"

// CreateSysJobReq 创建定时任务请求
type CreateSysJobReq struct {
	Name           string `json:"name" binding:"required"`
	JobGroup       string `json:"job_group" binding:"required"`
	CronExpression string `json:"cron_expression" binding:"required"`
	Command        string `json:"command" binding:"required"`
	Status         uint8  `json:"status" binding:"required,oneof=1 2"`
	Concurrent     bool   `json:"concurrent"`
	Description    string `json:"description"`
	RetryTimes    uint   `json:"retry_times"`
	RetryInterval uint   `json:"retry_interval"`
	Timeout       uint   `json:"timeout"`
}

// UpdateSysJobReq 更新定时任务请求
type UpdateSysJobReq struct {
	Id uint `json:"id"`
	CreateSysJobReq
}

// QuerySysJobReq 查询定时任务请求
type QuerySysJobReq struct {
	common.Pagination
	Name     string `json:"name" form:"name"`         // 任务名称
	JobGroup string `json:"job_group" form:"job_group"` // 任务分组
	Status   uint8  `json:"status" form:"status"`     // 状态
} 