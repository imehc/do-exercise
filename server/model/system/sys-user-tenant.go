package system

import (
	"time"

	"github.com/imehc/do-exercise/server/model"
	"gorm.io/gorm"
)

// SysUserTenant 用户-租户成员关系：一个用户当前在哪些租户里、以什么状态存在。
// 复合主键 (user_id, tenant_id)，同一用户在同一租户至多一条在册记录。
type SysUserTenant struct {
	UserId    string         `json:"user_id" gorm:"primaryKey;column:user_id;type:varchar(32);comment:用户ID"`
	TenantId  string         `json:"tenant_id" gorm:"primaryKey;column:tenant_id;type:varchar(32);comment:租户ID"`
	Status    bool           `json:"status" gorm:"column:status;default:true;comment:该租户内是否启用"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by;type:varchar(32);comment:创建人"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by;type:varchar(32);comment:更新人"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;comment:删除时间"`
}

func (*SysUserTenant) TableName() string {
	return "sys_user_tenant"
}

// BeforeCreate 回填操作人
func (u *SysUserTenant) BeforeCreate(tx *gorm.DB) (err error) {
	userId := model.CurrentUserID(tx)
	u.CreatedBy = userId
	u.UpdatedBy = userId
	return nil
}
