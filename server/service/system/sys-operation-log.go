package system

import (
	"context"
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
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
func (s *SysOperationLogService) Delete(id int) error {
	db := global.DB
	err := db.Delete(&system.SysOperationLog{}, id).Error
	if err != nil {
		global.Log.Error("删除操作日志失败", zap.Error(err))
		return err
	}
	return nil
}

// GetList 查询操作日志列表
func (s *SysOperationLogService) GetList(req common.Pagination) (*common.PageResult[system.SysOperationLog], error) {
	var logs []system.SysOperationLog
	var total int64
	db := global.DB.Model(&system.SysOperationLog{})
	db.Count(&total)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	db = db.
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Order("id ASC")
	err := db.Find(&logs).Error
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
