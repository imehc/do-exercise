package system

import (
	"context"
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysOperationLogService struct {
	context *context.Context
}

func (s *SysOperationLogService) WithContext(ctx context.Context) *SysOperationLogService {
	return &SysOperationLogService{
		context: &ctx,
	}
}

// Create 创建操作日志
func (s *SysOperationLogService) Create(req system.SysOperationLog) error {
	db := global.DB
	if s.context != nil {
		db = db.WithContext(*s.context)
	}

	err := db.Create(&req).Error
	if err != nil {
		global.Log.Error("创建操作日志失败", zap.Error(err))
		return err
	}
	return nil
}

// Delete 删除操作日志
func (s *SysOperationLogService) Delete(db *gorm.DB, id int) error {
	err := db.Delete(&system.SysOperationLog{}, id).Error
	if err != nil {
		global.Log.Error("删除操作日志失败", zap.Error(err))
		return err
	}
	return nil
}

// GetList 查询操作日志列表。
// 必须使用请求级 DB：sys_operation_log 在 tenantScopedTables 内，
// 租户插件会据此追加 tenant_id 过滤，租户只看到自己的日志；
// 平台超级管理员位于 platform 域，插件不追加条件，天然是全量视图。
// 此前这里直连 global.DB，语句不带请求上下文，隔离完全失效。
func (s *SysOperationLogService) GetList(db *gorm.DB, req common.Pagination) (*common.PageResult[system.SysOperationLog], error) {
	var logs []system.SysOperationLog
	var total int64

	// Count 用独立 builder，避免污染后续 Find 的状态
	countDB := db.Model(&system.SysOperationLog{})
	if err := countDB.Count(&total).Error; err != nil {
		return nil, errors.New("getFailed")
	}
	req.Normalize()
	err := db.Model(&system.SysOperationLog{}).
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Order("id DESC").
		Find(&logs).Error
	if err != nil {
		return nil, errors.New("getFailed")
	}

	return &common.PageResult[system.SysOperationLog]{
		Data: logs,
		Meta: common.PageMeta{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
		},
	}, nil
}
