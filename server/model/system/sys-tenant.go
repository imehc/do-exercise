package system

import (
	"time"

	"gorm.io/gorm"
)

// SysTenant 租户。租户目录表，全局共享，不参与行级租户隔离。
type SysTenant struct {
	TenantId   string         `json:"tenant_id" gorm:"primarykey;column:tenant_id;type:varchar(32);comment:租户ID"`
	Name       string         `json:"name" gorm:"not null;size:64;comment:租户名称"`
	Code       string         `json:"code" gorm:"not null;size:32;uniqueIndex:idx_sys_tenant_code,where:deleted_at IS NULL;comment:租户编码"`
	Status     bool           `json:"status" gorm:"default:true;comment:是否启用"`
	ExpireTime *time.Time     `json:"expire_time" gorm:"comment:过期时间"`
	Remark     string         `json:"remark" gorm:"size:255;default:null;comment:备注"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	CreatedBy  string         `json:"created_by" gorm:"column:created_by;type:varchar(32);comment:创建人"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	UpdatedBy  string         `json:"updated_by" gorm:"column:updated_by;type:varchar(32);comment:更新人"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index;comment:删除时间"`
}

func (*SysTenant) TableName() string {
	return "sys_tenant"
}

// IsTenantScoped 租户目录本身不做行级隔离
func (t *SysTenant) IsTenantScoped() bool {
	return false
}
