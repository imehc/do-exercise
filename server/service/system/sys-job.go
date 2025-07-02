package system

import (
	"errors"
	"fmt"

	"sync"

	"os/exec"

	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/common"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/model/system/request"
	"github.com/imehc/do-exercise/server/model/system/response"
	"github.com/imehc/do-exercise/server/util"
	"gorm.io/gorm"
)

var (
	scheduler     gocron.Scheduler
	schedulerOnce sync.Once
	jobMap        = make(map[uint]string) // jobId -> gocron.Job.ID().String()
	jobMapMutex   sync.Mutex
)

func getScheduler() gocron.Scheduler {
	schedulerOnce.Do(func() {
		var err error
		scheduler, err = gocron.NewScheduler()
		if err == nil {
			scheduler.Start()
		}
	})
	return scheduler
}

type SysJobService struct{}

// checkJobExist 检查任务是否存在
func (s *SysJobService) checkJobExist(db *gorm.DB, jobId uint) (*system.SysJob, error) {
	var job system.SysJob
	if err := db.First(&job, jobId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("job.notFound")
		}
		return nil, err
	}
	return &job, nil
}

// Create 创建定时任务
func (s *SysJobService) Create(req request.CreateSysJobReq) (*response.SysJobResp, error) {
	db := global.DB
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
	if err := db.Create(job).Error; err != nil {
		return nil, err
	}
	// 如果创建时状态为1，自动加入调度器
	if job.Status == 1 {
		gocronJob, err := s.addJobToScheduler(job)
		if err != nil {
			return nil, err
		}
		jobMapMutex.Lock()
		jobMap[job.Id] = gocronJob.ID().String()
		jobMapMutex.Unlock()
	}
	return s.convertToResp(job), nil
}

// Update 更新定时任务
func (s *SysJobService) Update(req request.UpdateSysJobReq) error {
	db := global.DB
	job, err := s.checkJobExist(db, req.Id)
	if err != nil {
		return err
	}
	// 先移除调度器中的旧任务（如有）
	jobMapMutex.Lock()
	if jobUUID, exists := jobMap[job.Id]; exists {
		_ = RemoveJobStr(getScheduler(), jobUUID)
		delete(jobMap, job.Id)
	}
	jobMapMutex.Unlock()
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
		return err
	}
	// 如新状态为1，重加到调度器
	if job.Status == 1 {
		gocronJob, err := s.addJobToScheduler(job)
		if err != nil {
			return err
		}
		jobMapMutex.Lock()
		jobMap[job.Id] = gocronJob.ID().String()
		jobMapMutex.Unlock()
	}
	return nil
}

// Delete 删除定时任务
func (s *SysJobService) Delete(id uint) error {
	db := global.DB
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}
	// 先从调度器移除
	jobMapMutex.Lock()
	if jobUUID, exists := jobMap[job.Id]; exists {
		_ = RemoveJobStr(getScheduler(), jobUUID)
		delete(jobMap, job.Id)
	}
	jobMapMutex.Unlock()
	// 再从数据库删除
	if err := db.Delete(&system.SysJob{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Get 获取单个定时任务
func (s *SysJobService) Get(id uint) (*response.SysJobResp, error) {
	db := global.DB
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return nil, err
	}

	return s.convertToResp(job), nil
}

// GetList 获取定时任务列表
func (s *SysJobService) GetList(req request.QuerySysJobReq) (*common.PageResult[response.SysJobResp], error) {
	db := global.DB
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
func (s *SysJobService) Start(id uint) error {
	db := global.DB
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}

	if job.Status == 1 {
		return errors.New("job.jobAlreadyStarted")
	}

	jobMapMutex.Lock()
	defer jobMapMutex.Unlock()
	if _, exists := jobMap[job.Id]; exists {
		return errors.New("job.alreadyInScheduler")
	}
	gocronJob, err := s.addJobToScheduler(job)
	if err != nil {
		return err
	}
	jobMap[job.Id] = gocronJob.ID().String()

	job.Status = 1
	if err := db.Save(job).Error; err != nil {
		return errors.New("job.startJobFailed")
	}

	return nil
}

// Stop 停止任务
func (s *SysJobService) Stop(id uint) error {
	db := global.DB
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}

	if job.Status == 2 {
		return errors.New("job.jobAlreadyStopped")
	}

	sch := getScheduler()
	jobMapMutex.Lock()
	if jobUUID, exists := jobMap[job.Id]; exists {
		_ = RemoveJobStr(sch, jobUUID)
		delete(jobMap, job.Id)
	}
	jobMapMutex.Unlock()

	job.Status = 2
	if err := db.Save(job).Error; err != nil {
		return errors.New("job.stopJobFailed")
	}

	return nil
}

// Execute 立即执行一次任务
func (s *SysJobService) Execute(id uint) error {
	db := global.DB
	job, err := s.checkJobExist(db, id)
	if err != nil {
		return err
	}
	return s.runJob(job)
}

// addJobToScheduler 添加任务到gocron调度器
func (s *SysJobService) addJobToScheduler(job *system.SysJob) (gocron.Job, error) {
	sch := getScheduler()
	gocronJob, err := sch.NewJob(
		gocron.CronJob(job.CronExpression, true), // 使用秒 即六位需要设置为true
		gocron.NewTask(func() { s.runJob(job) }),
		gocron.WithName(job.Name),
	)
	if err != nil {
		return nil, err
	}
	return gocronJob, nil
}

// runJob 执行任务命令
func (s *SysJobService) runJob(job *system.SysJob) error {
	db := global.DB

	// 记录上次执行时间
	now := time.Now()
	job.LastTime = &now
	job.Times++

	// 计算下次执行时间（如果有调度器信息）
	if scheduler != nil {
		jobMapMutex.Lock()
		if jobUUID, exists := jobMap[job.Id]; exists {
			uid := uuid.MustParse(jobUUID)
			for _, gocronJob := range scheduler.Jobs() {
				if gocronJob.ID() == uid {
					next, err := gocronJob.NextRun()
					if err == nil {
						nextUTC := next.UTC()
						job.NextTime = &nextUTC
					}
					break
				}
			}
		}
		jobMapMutex.Unlock()
	}

	db.Model(&system.SysJob{}).Where("id = ?", job.Id).Select("LastTime", "NextTime", "Times").Updates(job)

	if job.Command == "clean_empty_username_operation_logs" {
		return s.CleanEmptyUsernameOperationLogs()
	}
	// 支持执行shell命令
	cmd := exec.Command("/bin/sh", "-c", job.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("任务执行失败：", job.Name, job.Command, "错误：", err.Error(), "输出：", string(output))
		return err
	}
	fmt.Println("任务执行成功：", job.Name, job.Command, "输出：", string(output))
	return nil
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

// RemoveJobStr 用于通过字符串UUID移除任务
func RemoveJobStr(s gocron.Scheduler, id string) error {
	u, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.RemoveJob(u)
}

// 清理sys_operation_log表中username为空的任务实现
func (s *SysJobService) CleanEmptyUsernameOperationLogs() error {
	db := global.DB
	return db.Unscoped().Where("username = '' OR username IS NULL").Delete(&system.SysOperationLog{}).Error
}
