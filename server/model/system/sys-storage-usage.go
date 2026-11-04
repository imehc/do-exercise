package system

// ChargedBytes includes completed files and all outstanding upload reservations.
type SysStorageUsage struct {
	TenantId     string `json:"tenant_id" gorm:"primaryKey;size:32"`
	ChargedBytes int64  `json:"charged_bytes" gorm:"not null;default:0"`
}

func (*SysStorageUsage) TableName() string { return "sys_storage_usage" }
