package shared

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type jobHandler func(context.Context, *JobScheduler) error

var jobRegistry = map[string]jobHandler{
	"clean_empty_username_operation_logs": func(ctx context.Context, js *JobScheduler) error { return js.CleanEmptyUsernameOperationLogs(ctx) },
}

func IsRegisteredCommand(command string) bool { _, ok := jobRegistry[command]; return ok }
func RegisteredCommands() []string {
	result := make([]string, 0, len(jobRegistry))
	for key := range jobRegistry {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

var cronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

func ValidateJobSchedule(expression string) error { _, err := cronParser.Parse(expression); return err }

type JobScheduler struct {
	scheduler      gocron.Scheduler
	jobMap         map[uint]string
	revisions      map[uint]string
	jobMapMutex    sync.Mutex
	lifecycle      sync.Mutex
	stop           chan struct{}
	dispatcherDone chan struct{}
	workers        chan struct{}
	workReady      chan struct{}
	workerWG       sync.WaitGroup
}
type JobStats struct {
	LastTime time.Time `json:"last_time"`
	NextTime time.Time `json:"next_time"`
	Times    int64     `json:"times"`
	JobID    uint      `json:"job_id"`
}

var schedulerInstance *JobScheduler
var schedulerOnce sync.Once

func GetJobScheduler() *JobScheduler {
	schedulerOnce.Do(func() {
		scheduler, err := gocron.NewScheduler()
		if err != nil {
			panic(err)
		}
		schedulerInstance = &JobScheduler{scheduler: scheduler, jobMap: make(map[uint]string), revisions: make(map[uint]string), workers: make(chan struct{}, 8), workReady: make(chan struct{}, 1)}
		schedulerInstance.Start()
	})
	return schedulerInstance
}
func requireJobOwner(job *system.SysJob) error {
	if job.TenantId == "" {
		return errors.New("job tenant is required")
	}
	return nil
}
func jobContext(parent context.Context, job *system.SysJob) context.Context {
	ctx := context.WithValue(parent, global.ContextTenantIDKey, job.TenantId)
	if job.TenantId == global.PlatformTenantID {
		ctx = context.WithValue(ctx, global.ContextTenantBypassKey, true)
	}
	return context.WithValue(ctx, global.ContextUserIDKey, job.CreatedBy)
}
func platformMaintenanceDB(ctx context.Context) *gorm.DB {
	return util.TenantBypass(global.DB.WithContext(ctx))
}
func (js *JobScheduler) loadJobForMaintenance(id uint) (*system.SysJob, error) {
	var job system.SysJob
	err := platformMaintenanceDB(context.Background()).First(&job, id).Error
	return &job, err
}
func revision(job *system.SysJob) string {
	return fmt.Sprintf("%s/%d", job.CronExpression, job.UpdatedAt.UnixNano())
}
func (js *JobScheduler) AddJob(job *system.SysJob) error {
	if err := requireJobOwner(job); err != nil {
		return err
	}
	schedule, err := cronParser.Parse(job.CronExpression)
	if err != nil {
		return err
	}
	js.jobMapMutex.Lock()
	defer js.jobMapMutex.Unlock()
	if old, ok := js.jobMap[job.Id]; ok {
		if js.revisions[job.Id] == revision(job) {
			return nil
		}
		id, err := uuid.Parse(old)
		if err != nil {
			return err
		}
		if err = js.scheduler.RemoveJob(id); err != nil {
			return err
		}
		delete(js.jobMap, job.Id)
		delete(js.revisions, job.Id)
	}
	snapshot := *job
	scheduled, err := js.scheduler.NewJob(gocron.CronJob(job.CronExpression, true), gocron.NewTask(func() {
		// Only the cron occurrence near this invocation is eligible; delayed stale runs are skipped.
		now := time.Now()
		slot := schedule.Next(now.Add(-2 * time.Second))
		for next := schedule.Next(slot); !next.After(now); next = schedule.Next(slot) {
			slot = next
		}
		if slot.After(now) {
			return
		}
		if _, err := js.submit(&snapshot, &slot); err != nil && global.Log != nil {
			global.Log.Warn("定时任务未执行", zap.Uint("jobId", snapshot.Id), zap.Error(err))
		}
	}), gocron.WithName(job.Name))
	if err != nil {
		return err
	}
	js.jobMap[job.Id] = scheduled.ID().String()
	js.revisions[job.Id] = revision(job)
	return nil
}
func (js *JobScheduler) RemoveJob(id uint) error {
	js.jobMapMutex.Lock()
	defer js.jobMapMutex.Unlock()
	if raw, ok := js.jobMap[id]; ok {
		uid, err := uuid.Parse(raw)
		if err != nil {
			return err
		}
		if err = js.scheduler.RemoveJob(uid); err != nil {
			return err
		}
		delete(js.jobMap, id)
		delete(js.revisions, id)
	}
	return nil
}
func (js *JobScheduler) IsJobInScheduler(id uint) bool {
	js.jobMapMutex.Lock()
	defer js.jobMapMutex.Unlock()
	_, ok := js.jobMap[id]
	return ok
}
func (js *JobScheduler) GetNextRunTime(id uint) (*time.Time, error) {
	job, err := js.loadJobForMaintenance(id)
	if err != nil {
		return nil, err
	}
	if job.Status != 1 {
		return nil, errors.New("job stopped")
	}
	schedule, err := cronParser.Parse(job.CronExpression)
	if err != nil {
		return nil, err
	}
	next := schedule.Next(time.Now()).UTC()
	return &next, nil
}

func (js *JobScheduler) CleanEmptyUsernameOperationLogs(ctx context.Context) error {
	db, ok := ctx.Value(global.ContextDBKey).(*gorm.DB)
	if !ok {
		return errors.New("job execution transaction is required")
	}
	return db.Unscoped().Where("username = '' OR username IS NULL").Delete(&system.SysOperationLog{}).Error
}
func (js *JobScheduler) GetJobStats(id uint) (*JobStats, error) {
	job, err := js.loadJobForMaintenance(id)
	if err != nil {
		return nil, err
	}
	stats := &JobStats{JobID: id, Times: job.Times}
	if job.LastTime != nil {
		stats.LastTime = *job.LastTime
	}
	if next, err := js.GetNextRunTime(id); err == nil {
		stats.NextTime = *next
	}
	return stats, nil
}
func (js *JobScheduler) SyncOneJobStatsToDatabase(id uint) error {
	_, err := js.loadJobForMaintenance(id)
	return err
}
func (js *JobScheduler) SyncStatsToDatabase() error { return js.RestoreJobsFromDatabase() }
func (js *JobScheduler) RestoreJobsFromDatabase() error {
	var jobs []system.SysJob
	if err := platformMaintenanceDB(context.Background()).Where("status = ?", 1).Find(&jobs).Error; err != nil {
		return err
	}
	active := make(map[uint]bool)
	for i := range jobs {
		if err := util.RequireTenantUsable(global.DB, jobs[i].TenantId); err != nil {
			continue
		}
		active[jobs[i].Id] = true
		if err := js.AddJob(&jobs[i]); err != nil && global.Log != nil {
			global.Log.Warn("恢复任务失败", zap.Error(err))
		}
	}
	js.jobMapMutex.Lock()
	var removed []uint
	for id := range js.jobMap {
		if !active[id] {
			removed = append(removed, id)
		}
	}
	js.jobMapMutex.Unlock()
	for _, id := range removed {
		if err := js.RemoveJob(id); err != nil {
			return err
		}
	}
	return nil
}
func (js *JobScheduler) GetScheduler() gocron.Scheduler { return js.scheduler }
func (js *JobScheduler) Stop() {
	js.lifecycle.Lock()
	defer js.lifecycle.Unlock()
	if js.stop != nil {
		close(js.stop)
		<-js.dispatcherDone
		js.stop = nil
	}
	if js.scheduler != nil {
		js.scheduler.StopJobs()
		js.workerWG.Wait()
	}
}
func (js *JobScheduler) Start() {
	js.lifecycle.Lock()
	defer js.lifecycle.Unlock()
	if js.stop != nil {
		return
	}
	js.scheduler.Start()
	js.stop = make(chan struct{})
	stop := js.stop
	js.dispatcherDone = make(chan struct{})
	done := js.dispatcherDone
	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Second)
		var ticks uint
		defer ticker.Stop()
		for {
			select {
			case <-js.workReady:
				js.dispatchPending(stop)
			case <-ticker.C:
				ticks++
				js.dispatchPending(stop)
				if ticks%15 != 0 {
					continue
				}
				if err := js.RestoreJobsFromDatabase(); err != nil && global.Log != nil {
					global.Log.Warn("任务状态同步失败", zap.Error(err))
				}
			case <-stop:
				return
			}
		}
	}()
}
