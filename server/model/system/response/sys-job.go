package response

import "time"

// SysJobResp 定时任务响应
type SysJobResp struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	JobGroup       string     `json:"job_group"`
	CronExpression string     `json:"cron_expression"`
	Command        string     `json:"command"`
	Status         uint8      `json:"status"`
	Concurrent     bool       `json:"concurrent"`
	Description    string     `json:"description"`
	LastTime       *time.Time `json:"last_time,omitzero"`
	NextTime       *time.Time `json:"next_time,omitzero"`
	Times          int64      `json:"times"`
	RetryTimes     uint       `json:"retry_times"`
	RetryInterval  uint       `json:"retry_interval"`
	Timeout        uint       `json:"timeout"`
	CreatedAt      time.Time  `json:"created_at"`
	CreatedBy      string     `json:"created_by"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UpdatedBy      string     `json:"updated_by"`
}
