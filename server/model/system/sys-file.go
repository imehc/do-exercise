package system

import "time"

// SysFile records the server-assigned owner and upload limits before issuing a signature.
type SysFile struct {
	ObjectKey       string     `gorm:"primaryKey;size:512"`
	TenantId        string     `gorm:"not null;size:32;index"`
	UserId          string     `gorm:"not null;size:32;index"`
	FileName        string     `gorm:"not null;size:255"`
	ContentType     string     `gorm:"not null;size:128"`
	MaxSize         int64      `gorm:"not null"`
	Public          bool       `gorm:"not null;default:false"`
	Size            int64      `gorm:"not null;default:0"`
	State           string     `gorm:"not null;size:16;default:legacy;index:idx_file_gc,priority:1"`
	UploadExpiresAt *time.Time `gorm:"index"`
	DeleteAfter     *time.Time `gorm:"index:idx_file_gc,priority:2"`
	CompletedAt     *time.Time
	DeletedAt       *time.Time
	CreatedAt       time.Time
}

func (*SysFile) TableName() string { return "sys_file" }
