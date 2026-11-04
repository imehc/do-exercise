package system

import (
	"errors"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
)

type SysJobService struct{}

// checkJobExist 检查任务是否存在
func (s *SysJobService) checkJobExist(db *gorm.DB, jobId uint) (*system.SysJob, error) {
	var job system.SysJob
	if err := db.First(&job, jobId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job.notFound")
		}
		return nil, errors.New("job.getJobFailed")
	}
	return &job, nil
}

// Create 创建定时任务
func (s *SysJobService) Create(db *gorm.DB, req request.CreateSysJobReq) (*response.SysJobResp, error) {
	if err := shared.ValidateExecutionOptions(req.RetryTimes, req.RetryInterval, req.Timeout); err != nil {
		return nil, err
	}
	// 命令必须在白名单内。放行任意字符串等同于开放容器内的 shell。
	if err := shared.ValidateJobSchedule(req.CronExpression); err != nil {
		return nil, err
	}
	if !shared.IsRegisteredCommand(req.Command) {
		return nil, errors.New("job.unregisteredCommand")
	}
	job := &system.SysJob{
		Name:           req.Name,
		JobGroup:       req.JobGroup,
		CronExpression: req.CronExpression,
		Command:        req.Command,
		Status:         req.Status,
		Concurrent:     req.Concurrent,
		Description:    req.Description,
		RetryTimes:     req.RetryTimes,
		RetryInterval:  req.RetryInterval,
		Timeout:        req.Timeout,
	}
	// 归属显式化。租户上下文由插件回填；平台超管没有租户上下文，回填不出值，
	// 这里显式落成平台租户——平台维护类任务（日志清理等）需要跨租户执行，
	// 空 tenant_id 则是谁都认领不了的孤儿行，调度器会直接拒绝。
	job.TenantId = model.CurrentTenantID(db)
	if job.TenantId == "" && model.CurrentIsSuperAdmin(db) {
		job.TenantId = global.PlatformTenantID
	}
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := enforceJobQuota(tx, job.TenantId, 1); err != nil {
			return err
		}
		return tx.Create(job).Error
	}); err != nil {
		return nil, err
	}

	// 如果创建时状态为1，自动加入调度器
	if job.Status == 1 {
		jobScheduler := shared.GetJobScheduler()
		if err := jobScheduler.AddJob(job); err != nil {
			return nil, errors.New("job.createJobFailed")
		}
	}
	return s.convertToResp(job), nil
}

// Update 更新定时任务
func (s *SysJobService) Update(db *gorm.DB, id uint, req request.UpdateSysJobReq) error {
	if err := shared.ValidateExecutionOptions(req.RetryTimes, req.RetryInterval, req.Timeout); err != nil {
		return err
	}
	// 命令必须在白名单内，防止通过更新绕过创建时的校验
	if err := shared.ValidateJobSchedule(req.CronExpression); err != nil {
		return err
	}
	if !shared.IsRegisteredCommand(req.Command) {
		return errors.New("job.unregisteredCommand")
	}
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}
	// 先移除调度器中的旧任务（如有）
	jobScheduler := shared.GetJobScheduler()
	_ = jobScheduler.RemoveJob(job.Id)

	// 更新字段
	job.Name = req.Name
	job.JobGroup = req.JobGroup
	job.CronExpression = req.CronExpression
	job.Command = req.Command
	job.Status = req.Status
	job.Concurrent = req.Concurrent
	job.Description = req.Description
	job.RetryTimes = req.RetryTimes
	job.RetryInterval = req.RetryInterval
	job.Timeout = req.Timeout
	if err := db.Save(job).Error; err != nil {
		return errors.New("job.updateJobFailed")
	}
	// 如新状态为1，重加到调度器
	if job.Status == 1 {
		if err := jobScheduler.AddJob(job); err != nil {
			return errors.New("job.updateJobFailed")
		}
	}
	return nil
}

// Delete 删除定时任务
func (s *SysJobService) Delete(db *gorm.DB, id uint) error {
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}
	// 先从调度器移除
	jobScheduler := shared.GetJobScheduler()
	_ = jobScheduler.RemoveJob(job.Id)

	// 再从数据库删除
	if err := db.Delete(&system.SysJob{}, id).Error; err != nil {
		return errors.New("job.deleteJobFailed")
	}
	return nil
}

// Get 获取单个定时任务
func (s *SysJobService) Get(db *gorm.DB, id uint) (*response.SysJobResp, error) {
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return nil, err
	}

	resp := s.convertToResp(job)

	// 从数据库获取已提交的执行统计
	jobScheduler := shared.GetJobScheduler()
	if stats, err := jobScheduler.GetJobStats(job.Id); err == nil {
		// 所有实例共享数据库统计
		if !stats.LastTime.IsZero() {
			resp.LastTime = &stats.LastTime
		}
		if !stats.NextTime.IsZero() {
			resp.NextTime = &stats.NextTime
		}
		resp.Times = stats.Times
	}

	return resp, nil
}

// GetList 获取定时任务列表
func (s *SysJobService) GetList(db *gorm.DB, req request.QuerySysJobReq) (*common.PageResult[response.SysJobResp], error) {
	var jobs []system.SysJob
	var total int64

	// 构建查询条件
	query := db.Model(&system.SysJob{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.JobGroup != "" {
		query = query.Where("job_group = ?", req.JobGroup)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("job.getJobListFailed")
	}

	req.Normalize()
	// 获取分页数据
	if err := query.
		Scopes(util.Paginate(req.PageSize, req.Page)).
		Find(&jobs).Error; err != nil {
		return nil, errors.New("job.getJobListFailed")
	}

	// 转换为响应结构
	data := make([]response.SysJobResp, len(jobs))
	for i, job := range jobs {
		data[i] = *s.convertToResp(&job)
	}

	return &common.PageResult[response.SysJobResp]{
		Data: data,
		Meta: common.PageMeta{
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}, nil
}

// Start 启动任务
func (s *SysJobService) Start(db *gorm.DB, id uint) error {
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}

	if job.Status == 1 {
		return errors.New("job.jobAlreadyStarted")
	}

	jobScheduler := shared.GetJobScheduler()
	job.Status = 1
	if err := db.Save(job).Error; err != nil {
		return errors.New("job.startJobFailed")
	}
	if err := jobScheduler.AddJob(job); err != nil {
		return errors.New("job.startJobFailed")
	}

	return nil
}

// Stop 停止任务
func (s *SysJobService) Stop(db *gorm.DB, id uint) error {
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}

	if job.Status == 2 {
		return errors.New("job.jobAlreadyStopped")
	}

	jobScheduler := shared.GetJobScheduler()

	// 1. 直接同步 Redis 统计到数据库（只读，不写入Redis）
	err = jobScheduler.SyncOneJobStatsToDatabase(id)
	if err != nil {
		return errors.New("job.stopJobFailed")
	}
	// 2. 再移除调度器任务
	_ = jobScheduler.RemoveJob(id)

	if err := db.Model(job).
		Update("Status", 2).
		Update("NextTime", nil).
		Error; err != nil {
		return errors.New("job.stopJobFailed")
	}

	return nil
}

// Execute 立即执行一次任务
func (s *SysJobService) Execute(db *gorm.DB, id uint) error {
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}

	_, err = shared.GetJobScheduler().SubmitJob(job)
	return err
}

// GetJobStats 获取任务统计信息
func (s *SysJobService) GetJobStats(db *gorm.DB, id uint) (*shared.JobStats, error) {
	// 检查任务是否存在
	_, err := s.checkJobExist(db, id)
	if err != nil {
		return nil, err
	}

	// 从Redis获取统计信息
	jobScheduler := shared.GetJobScheduler()
	stats, err := jobScheduler.GetJobStats(id)
	if err != nil {
		return nil, errors.New("job.getJobFailed")
	}
	return stats, nil
}

// SyncStatsToDatabase 手动同步统计信息到数据库
func (s *SysJobService) SyncStatsToDatabase() error {
	jobScheduler := shared.GetJobScheduler()
	return jobScheduler.SyncStatsToDatabase()
}

// convertToResp 将数据库模型转换为响应结构
func (s *SysJobService) convertToResp(job *system.SysJob) *response.SysJobResp {
	return &response.SysJobResp{
		Id:             job.Id,
		Name:           job.Name,
		JobGroup:       job.JobGroup,
		CronExpression: job.CronExpression,
		Command:        job.Command,
		Status:         job.Status,
		Concurrent:     job.Concurrent,
		Description:    job.Description,
		LastTime:       job.LastTime,
		NextTime:       job.NextTime,
		Times:          job.Times,
		RetryTimes:     job.RetryTimes,
		RetryInterval:  job.RetryInterval,
		Timeout:        job.Timeout,
		CreatedAt:      job.CreatedAt,
		CreatedBy:      job.CreatedBy,
		UpdatedAt:      job.UpdatedAt,
		UpdatedBy:      job.UpdatedBy,
	}
}

func (s *SysJobService) GetRuns(db *gorm.DB, id uint, req common.Pagination) (*common.PageResult[system.SysJobRun], error) {
	if _, err := s.checkJobExist(db, id); err != nil {
		return nil, err
	}
	req.Normalize()
	var total int64
	var runs []system.SysJobRun
	if err := db.Model(&system.SysJobRun{}).Where("job_id = ?", id).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := db.Where("job_id = ?", id).Order("created_at DESC,id DESC").Scopes(util.Paginate(req.PageSize, req.Page)).Find(&runs).Error; err != nil {
		return nil, err
	}
	return &common.PageResult[system.SysJobRun]{Data: runs, Meta: common.PageMeta{Total: total, Page: req.Page, PageSize: req.PageSize}}, nil
}
