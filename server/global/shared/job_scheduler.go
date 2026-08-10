package shared

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// jobHandler 定时任务处理函数。ctx 用于传递超时，处理函数应当尊重其取消信号。
type jobHandler func(ctx context.Context, js *JobScheduler) error

// jobRegistry 是允许被调度执行的任务白名单。
//
// SysJob.Command 只能取本表中的 key。早期实现把该字段直接交给 `/bin/sh -c`，
// 于是任何能创建定时任务的账号都等价于拿到容器内的 shell，
// 且能挂在 cron 上周期执行——管理员账号被盗、后台 XSS、CSRF 任一路径都会升级为 RCE。
// 新增任务请在此登记，不要恢复通用 shell 执行。
var jobRegistry = map[string]jobHandler{
	"clean_empty_username_operation_logs": func(ctx context.Context, js *JobScheduler) error {
		return js.CleanEmptyUsernameOperationLogs(ctx)
	},
}

// IsRegisteredCommand 判断任务命令是否在白名单内，供服务层在创建/更新时校验
func IsRegisteredCommand(command string) bool {
	_, ok := jobRegistry[command]
	return ok
}

// RegisteredCommands 返回全部可执行的任务名，用于错误提示与前端选项
func RegisteredCommands() []string {
	commands := make([]string, 0, len(jobRegistry))
	for name := range jobRegistry {
		commands = append(commands, name)
	}
	sort.Strings(commands)
	return commands
}

// JobScheduler 定时任务调度器
type JobScheduler struct {
	scheduler     gocron.Scheduler
	jobMap        map[uint]string // jobId -> gocron.Job.ID().String()
	jobMapMutex   sync.Mutex
	syncTicker    *time.Ticker
	syncStop      chan struct{} // 关闭后通知同步 goroutine 退出
	syncRunning   bool          // 防止 Start() 重复启动 ticker 造成泄漏
	syncTickerMux sync.Mutex
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
		}

		// 启动定时同步任务
		schedulerInstance.startSyncTicker()
	})
	return schedulerInstance
}

// startSyncTicker 启动定时同步任务。
// 用 syncRunning 标志防重入：Start() 反复调用不会反复新建 ticker。
func (js *JobScheduler) startSyncTicker() {
	js.syncTickerMux.Lock()
	defer js.syncTickerMux.Unlock()
	if js.syncRunning {
		return
	}
	js.syncRunning = true
	js.syncStop = make(chan struct{})
	js.syncTicker = time.NewTicker(1 * time.Minute) // 每分钟同步一次
	go func() {
		for {
			select {
			case <-js.syncTicker.C:
				js.SyncStatsToDatabase()
			case <-js.syncStop:
				js.syncTicker.Stop()
				return
			}
		}
	}()
}

// stopSyncTicker 停止定时同步任务。
// 关闭 channel 而非无缓冲发送，避免接收方已退出时永久阻塞。
func (js *JobScheduler) stopSyncTicker() {
	js.syncTickerMux.Lock()
	defer js.syncTickerMux.Unlock()
	if !js.syncRunning {
		return
	}
	js.syncRunning = false
	close(js.syncStop)
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

// updateJobStatsToDB 封装原子累加并更新结构体。
// 用 SQL 表达式 `times = times + ?` 而非读-改-写，避免并发同步丢次数。
func (js *JobScheduler) updateJobStatsToDB(db *gorm.DB, jobId uint, stats *JobStats) error {
	return db.Model(&system.SysJob{}).
		Where("id = ?", jobId).
		Updates(map[string]interface{}{
			"last_time": stats.LastTime,
			"next_time": stats.NextTime,
			"times":     gorm.Expr("times + ?", stats.Times),
		}).Error
}

// SyncStatsToDatabase 同步统计信息到数据库
func (js *JobScheduler) SyncStatsToDatabase() error {
	db := global.DB.
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
			global.Log.Error("同步任务统计信息失败",
				zap.Uint("jobId", jobId),
				zap.Error(err))
			continue
		}
		if err := js.updateJobStats(jobId, stats.LastTime, stats.NextTime, 0); err != nil {
			global.Log.Error("重置任务统计计数失败",
				zap.Uint("jobId", jobId),
				zap.Error(err))
		}
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
	if err := js.updateJobStats(job.Id, now, nextTime, 0); err != nil {
		global.Log.Error("初始化任务统计信息失败",
			zap.Uint("jobId", job.Id),
			zap.Error(err))
	}

	return nil
}

// RemoveJob 从调度器中移除任务
func (js *JobScheduler) RemoveJob(jobId uint) error {
	js.jobMapMutex.Lock()
	defer js.jobMapMutex.Unlock()

	if jobUUID, exists := js.jobMap[jobId]; exists {
		u, err := uuid.Parse(jobUUID) // 不再用 uuid.MustParse，畸形值不 panic
		if err != nil {
			global.Log.Error("解析任务UUID失败",
				zap.Uint("jobId", jobId),
				zap.Error(err))
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
		uid, err := uuid.Parse(jobUUID) // 不再用 uuid.MustParse，畸形值不 panic
		if err != nil {
			return nil, err
		}
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
		if err := js.updateJobStats(job.Id, now, nextTime, 1); err != nil {
			global.Log.Error("初始化任务执行统计失败",
				zap.Uint("jobId", job.Id),
				zap.Error(err))
		}
	} else {
		// 累加执行次数
		if err := js.updateJobStats(job.Id, now, nextTime, stats.Times+1); err != nil {
			global.Log.Error("累加任务执行统计失败",
				zap.Uint("jobId", job.Id),
				zap.Error(err))
		}
	}

	// 执行任务：只允许白名单内已登记的处理函数，不再解释执行任意 shell 命令
	handler, ok := jobRegistry[job.Command]
	if !ok {
		global.Log.Error("任务命令未登记，拒绝执行",
			zap.String("job", job.Name),
			zap.String("command", job.Command))
		return fmt.Errorf("unregistered job command: %s", job.Command)
	}

	// 当设置了Timeout时在上下文中施加超时，避免任务永久卡死；Timeout为0保持原有行为
	ctx := context.Background()
	cancel := func() {}
	if job.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(job.Timeout)*time.Second)
	}
	defer cancel()

	if err := handler(ctx, js); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			global.Log.Error("任务执行超时",
				zap.String("job", job.Name),
				zap.String("command", job.Command))
		} else {
			global.Log.Error("任务执行失败",
				zap.String("job", job.Name),
				zap.String("command", job.Command),
				zap.Error(err))
		}
		return err
	}

	global.Log.Info("任务执行成功",
		zap.String("job", job.Name),
		zap.String("command", job.Command))
	return nil
}

// CleanEmptyUsernameOperationLogs 清理sys_operation_log表中username为空的任务实现
func (js *JobScheduler) CleanEmptyUsernameOperationLogs(ctx context.Context) error {
	db := global.DB.WithContext(ctx)
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
			global.Log.Error("恢复任务失败",
				zap.String("job", job.Name),
				zap.Error(err))
			continue
		}
		global.Log.Info("成功恢复任务", zap.String("job", job.Name))
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
