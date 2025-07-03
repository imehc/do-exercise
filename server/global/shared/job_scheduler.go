package shared

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"gorm.io/gorm"
)

// JobScheduler 定时任务调度器
type JobScheduler struct {
	scheduler   gocron.Scheduler
	jobMap      map[uint]string // jobId -> gocron.Job.ID().String()
	jobMapMutex sync.Mutex
	syncTicker  *time.Ticker
	stopChan    chan bool
}

// JobStats 任务统计信息
type JobStats struct {
	LastTime time.Time `json:"last_time"`
	NextTime time.Time `json:"next_time"`
	Times    int64     `json:"times"`
	JobID    uint      `json:"job_id"`
}

var (
	schedulerInstance *JobScheduler
	schedulerOnce     sync.Once
)

// GetJobScheduler 获取定时任务调度器单例
func GetJobScheduler() *JobScheduler {
	schedulerOnce.Do(func() {
		scheduler, err := gocron.NewScheduler()
		if err == nil {
			scheduler.Start()
		}
		schedulerInstance = &JobScheduler{
			scheduler: scheduler,
			jobMap:    make(map[uint]string),
			stopChan:  make(chan bool),
		}

		// 启动定时同步任务
		schedulerInstance.startSyncTicker()
	})
	return schedulerInstance
}

// startSyncTicker 启动定时同步任务
func (js *JobScheduler) startSyncTicker() {
	js.syncTicker = time.NewTicker(1 * time.Minute) // 每分钟同步一次
	go func() {
		for {
			select {
			case <-js.syncTicker.C:
				js.SyncStatsToDatabase()
			case <-js.stopChan:
				js.syncTicker.Stop()
				return
			}
		}
	}()
}

// stopSyncTicker 停止定时同步任务
func (js *JobScheduler) stopSyncTicker() {
	if js.syncTicker != nil {
		js.stopChan <- true
	}
}

// getRedisKey 获取Redis键名
func (js *JobScheduler) getRedisKey(jobId uint) string {
	return fmt.Sprintf("job_stats:%d", jobId)
}

// updateJobStats 更新任务统计信息到Redis
func (js *JobScheduler) updateJobStats(jobId uint, lastTime time.Time, nextTime time.Time, times int64) error {
	ctx := context.Background()
	key := js.getRedisKey(jobId)

	// 统一为本地时区
	localLast := lastTime.Local()
	localNext := nextTime.Local()

	pipe := global.Redis.Pipeline()
	pipe.HSet(ctx, key, "last_time", localLast.Format(time.RFC3339))
	pipe.HSet(ctx, key, "next_time", localNext.Format(time.RFC3339))
	pipe.HSet(ctx, key, "times", times)
	pipe.HSet(ctx, key, "job_id", jobId)
	pipe.Expire(ctx, key, 24*time.Hour) // 设置24小时过期时间

	_, err := pipe.Exec(ctx)
	return err
}

// getJobStats 从Redis获取任务统计信息
func (js *JobScheduler) getJobStats(jobId uint) (*JobStats, error) {
	ctx := context.Background()
	key := js.getRedisKey(jobId)

	result, err := global.Redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("job stats not found")
	}

	lastTime, _ := time.Parse(time.RFC3339, result["last_time"])
	nextTime, _ := time.Parse(time.RFC3339, result["next_time"])
	times, _ := strconv.ParseInt(result["times"], 10, 64)
	jobID, _ := strconv.ParseUint(result["job_id"], 10, 32)

	return &JobStats{
		LastTime: lastTime,
		NextTime: nextTime,
		Times:    times,
		JobID:    uint(jobID),
	}, nil
}

// updateJobStatsToDB 封装累加并结构体更新逻辑
func (js *JobScheduler) updateJobStatsToDB(db *gorm.DB, jobId uint, stats *JobStats) error {
	var job system.SysJob
	if err := db.Where("id = ?", jobId).First(&job).Error; err != nil {
		return err
	}
	newTimes := job.Times + stats.Times
	return db.Model(&system.SysJob{}).
		Where("id = ?", jobId).
		Select("LastTime", "NextTime", "Times").
		Updates(&system.SysJob{
			LastTime: &stats.LastTime,
			NextTime: &stats.NextTime,
			Times:    newTimes,
		}).Error
}

// SyncStatsToDatabase 同步统计信息到数据库
func (js *JobScheduler) SyncStatsToDatabase() error {
	db := global.DB.
		// Session(&gorm.Session{
		// 	Logger: logger.Default.LogMode(logger.Info),
		// }).
		WithContext(context.Background())
	js.jobMapMutex.Lock()
	jobIds := make([]uint, 0, len(js.jobMap))
	for jobId := range js.jobMap {
		jobIds = append(jobIds, jobId)
	}
	js.jobMapMutex.Unlock()
	for _, jobId := range jobIds {
		stats, err := js.getJobStats(jobId)
		if err != nil {
			continue // 跳过不存在的统计信息
		}
		if err := js.updateJobStatsToDB(db, jobId, stats); err != nil {
			fmt.Printf("同步任务 %d 统计信息失败: %v\n", jobId, err)
			continue
		}
		js.updateJobStats(jobId, stats.LastTime, stats.NextTime, 0)
		fmt.Printf("成功同步任务 %d 统计信息到数据库\n", jobId)
	}
	return nil
}

// AddJob 添加任务到调度器
func (js *JobScheduler) AddJob(job *system.SysJob) error {
	gocronJob, err := js.scheduler.NewJob(
		gocron.CronJob(job.CronExpression, true), // 使用秒 即六位需要设置为true
		gocron.NewTask(func() { js.runJob(job) }),
		gocron.WithName(job.Name),
	)
	if err != nil {
		return err
	}

	js.jobMapMutex.Lock()
	js.jobMap[job.Id] = gocronJob.ID().String()
	js.jobMapMutex.Unlock()

	// 初始化Redis统计信息
	now := time.Now()
	nextTime, _ := gocronJob.NextRun()
	js.updateJobStats(job.Id, now, nextTime, 0)

	return nil
}

// RemoveJob 从调度器中移除任务
func (js *JobScheduler) RemoveJob(jobId uint) error {
	js.jobMapMutex.Lock()
	defer js.jobMapMutex.Unlock()

	if jobUUID, exists := js.jobMap[jobId]; exists {
		u, err := uuid.Parse(jobUUID)
		if err != nil {
			return err
		}
		err = js.scheduler.RemoveJob(u)
		if err == nil {
			delete(js.jobMap, jobId)
			// 清理Redis统计信息
			ctx := context.Background()
			global.Redis.Del(ctx, js.getRedisKey(jobId))
		}
		return err
	}
	return nil
}

// IsJobInScheduler 检查任务是否在调度器中
func (js *JobScheduler) IsJobInScheduler(jobId uint) bool {
	js.jobMapMutex.Lock()
	defer js.jobMapMutex.Unlock()
	_, exists := js.jobMap[jobId]
	return exists
}

// GetNextRunTime 获取任务下次执行时间
func (js *JobScheduler) GetNextRunTime(jobId uint) (*time.Time, error) {
	js.jobMapMutex.Lock()
	defer js.jobMapMutex.Unlock()

	if jobUUID, exists := js.jobMap[jobId]; exists {
		uid := uuid.MustParse(jobUUID)
		for _, gocronJob := range js.scheduler.Jobs() {
			if gocronJob.ID() == uid {
				next, err := gocronJob.NextRun()
				if err != nil {
					return nil, err
				}
				nextUTC := next.UTC()
				return &nextUTC, nil
			}
		}
	}
	return nil, errors.New("job not found in scheduler")
}

// GetJobStatsFromRedis 获取任务统计信息（从Redis）
func (js *JobScheduler) GetJobStatsFromRedis(jobId uint) (*JobStats, error) {
	return js.getJobStats(jobId)
}

// runJob 执行任务命令
func (js *JobScheduler) runJob(job *system.SysJob) error {
	// 记录执行开始时间
	now := time.Now()

	// 获取下次执行时间
	var nextTime time.Time
	if next, err := js.GetNextRunTime(job.Id); err == nil {
		nextTime = *next
	} else {
		nextTime = now.Add(time.Hour) // 默认1小时后
	}

	// 更新Redis统计信息
	stats, err := js.getJobStats(job.Id)
	if err != nil {
		// 如果Redis中没有统计信息，创建新的
		js.updateJobStats(job.Id, now, nextTime, 1)
	} else {
		// 累加执行次数
		js.updateJobStats(job.Id, now, nextTime, stats.Times+1)
	}

	// 执行任务
	if job.Command == "clean_empty_username_operation_logs" {
		return js.CleanEmptyUsernameOperationLogs()
	}

	// 支持执行shell命令
	cmd := exec.Command("/bin/sh", "-c", job.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("任务执行失败：%s, 命令：%s, 错误：%s, 输出：%s\n",
			job.Name, job.Command, err.Error(), string(output))
		return err
	}
	fmt.Printf("任务执行成功：%s, 命令：%s, 输出：%s\n",
		job.Name, job.Command, string(output))
	return nil
}

// CleanEmptyUsernameOperationLogs 清理sys_operation_log表中username为空的任务实现
func (js *JobScheduler) CleanEmptyUsernameOperationLogs() error {
	db := global.DB
	return db.Unscoped().Where("username = '' OR username IS NULL").Delete(&system.SysOperationLog{}).Error
}

// ExecuteJob 立即执行一次任务
func (js *JobScheduler) ExecuteJob(job *system.SysJob) error {
	return js.runJob(job)
}

// GetScheduler 获取底层调度器（用于兼容性）
func (js *JobScheduler) GetScheduler() gocron.Scheduler {
	return js.scheduler
}

// RestoreJobsFromDatabase 从数据库恢复定时任务
func (js *JobScheduler) RestoreJobsFromDatabase() error {
	db := global.DB
	var jobs []system.SysJob

	// 查询状态为1（正常）的任务
	if err := db.Where("status = ?", 1).Find(&jobs).Error; err != nil {
		return err
	}

	// 将任务添加到调度器
	for _, job := range jobs {
		if err := js.AddJob(&job); err != nil {
			fmt.Printf("恢复任务失败: %s, 错误: %v\n", job.Name, err)
			continue
		}
		fmt.Printf("成功恢复任务: %s\n", job.Name)
	}

	return nil
}

// Stop 停止调度器
func (js *JobScheduler) Stop() {
	js.stopSyncTicker()
	if js.scheduler != nil {
		js.scheduler.Shutdown()
	}
}

// Start 启动调度器
func (js *JobScheduler) Start() {
	if js.scheduler != nil {
		js.scheduler.Start()
	}
	js.startSyncTicker()
}

// SyncOneJobStatsToDatabase 同步单个任务统计信息到数据库
func (js *JobScheduler) SyncOneJobStatsToDatabase(jobId uint) error {
	db := global.DB.
		// Session(&gorm.Session{
		// 	Logger: logger.Default.LogMode(logger.Info),
		// }).
		WithContext(context.Background())
	stats, err := js.getJobStats(jobId)
	if err != nil {
		return err
	}
	if err := js.updateJobStatsToDB(db, jobId, stats); err != nil {
		return err
	}
	// 同步成功后删除Redis统计
	ctx := context.Background()
	global.Redis.Del(ctx, js.getRedisKey(jobId))
	return nil
}

// UpdateJobStats 更新任务统计信息到Redis（公有方法）
func (js *JobScheduler) UpdateJobStats(jobId uint, lastTime time.Time, nextTime time.Time, times int64) error {
	return js.updateJobStats(jobId, lastTime, nextTime, times)
}
