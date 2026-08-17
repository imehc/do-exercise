package model

import (
	"strings"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm"
)

// CurrentUserID 从 GORM 事务上下文安全地取出当前操作人 ID。
// 上下文可能没有值，或值不是 string（迁移、种子、定时任务等不经
// ContextMiddleware 的写入路径），裸断言会 panic，这里做 ok 检查并兜底为空串。
func CurrentUserID(tx *gorm.DB) string {
	v := tx.Statement.Context.Value(global.ContextUserIDKey)
	userId, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(userId)
}

// CurrentTenantID 从 GORM 事务上下文安全地取出当前租户 ID。
// 与 CurrentUserID 相同的兜底策略；不经请求上下文的路径返回空串。
func CurrentTenantID(tx *gorm.DB) string {
	v := tx.Statement.Context.Value(global.ContextTenantIDKey)
	tenantId, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(tenantId)
}

type IdWrapper struct {
	Id uint `json:"id" gorm:"primarykey;comment:主键"`
}

type ControlWrapper struct {
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime;comment:创建时间"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by;type:varchar(32);comment:创建人"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;comment:更新时间"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by;type:varchar(32);comment:更新人"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index;comment:删除时间"`
	DeletedBy string         `json:"deleted_by" gorm:"column:deleted_by;type:varchar(32);comment:删除人"`
}
