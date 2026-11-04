package system

import "time"

// SysJobRun is a durable occurrence. Database handlers and its terminal state
// commit in one transaction; a crashed attempt leaves no committed effects.
type SysJobRun struct {
	ID            string     `json:"id" gorm:"primaryKey;size:32"`
	JobID         uint       `json:"job_id" gorm:"not null;index;uniqueIndex:idx_job_occurrence,priority:1"`
	TenantId      string     `json:"tenant_id" gorm:"not null;size:32;index"`
	Occurrence    string     `json:"-" gorm:"not null;size:64;uniqueIndex:idx_job_occurrence,priority:2"`
	Scheduled     bool       `json:"scheduled" gorm:"not null"`
	Status        string     `json:"status" gorm:"not null;size:16;index:idx_job_run_due,priority:1"`
	Attempts      uint       `json:"attempts" gorm:"not null;default:0"`
	Command       string     `json:"-" gorm:"not null;size:128"`
	RetryTimes    uint       `json:"-" gorm:"not null"`
	RetryInterval uint       `json:"-" gorm:"not null"`
	Timeout       uint       `json:"-" gorm:"not null"`
	NextAttemptAt time.Time  `json:"next_attempt_at" gorm:"not null;index:idx_job_run_due,priority:2"`
	LastError     string     `json:"last_error" gorm:"not null;default:'';size:1024"`
	CreatedAt     time.Time  `json:"created_at"`
	FinishedAt    *time.Time `json:"finished_at,omitempty" gorm:"index"`
}

func (*SysJobRun) TableName() string { return "sys_job_run" }
