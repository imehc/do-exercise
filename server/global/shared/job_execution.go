package shared

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const activeRuns = "status IN ('queued','retrying')"

func ValidateExecutionOptions(retries, interval, timeout uint) error {
	if retries > 5 || interval > 300 || timeout > 3600 {
		return errors.New("job.invalidExecutionOptions")
	}
	return nil
}

// Admission is short and serialized per job, while independent occurrences may
// run concurrently. Overlapping occurrences are skipped instead of accumulating.
func (js *JobScheduler) submit(snapshot *system.SysJob, slot *time.Time) (*system.SysJobRun, error) {
	if err := requireJobOwner(snapshot); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var run *system.SysJobRun
	err := platformMaintenanceDB(ctx).Transaction(func(tx *gorm.DB) error {
		var job system.SysJob
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&job, snapshot.Id).Error; err != nil {
			return err
		}
		if job.TenantId != snapshot.TenantId {
			return errors.New("job ownership changed")
		}
		if slot != nil {
			if job.Status != 1 || revision(&job) != revision(snapshot) {
				return errors.New("job schedule changed or stopped")
			}
			if job.ScheduledAt != nil && !job.ScheduledAt.Before(*slot) {
				return nil
			}
		}
		if err := util.RequireTenantUsable(tx, job.TenantId); err != nil {
			return err
		}
		if !IsRegisteredCommand(job.Command) {
			return errors.New("unregistered job command")
		}
		if err := ValidateExecutionOptions(job.RetryTimes, job.RetryInterval, job.Timeout); err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&system.SysJobRun{}).Where("job_id = ?", job.Id).Where(activeRuns).Count(&count).Error; err != nil {
			return err
		}
		limit := int64(1)
		if job.Concurrent {
			limit = 4
		}
		id := util.NextID()
		occurrence := id
		if slot != nil {
			occurrence = strconv.FormatInt(slot.UnixNano(), 10)
		}
		now := time.Now().UTC()
		run = &system.SysJobRun{ID: id, JobID: job.Id, TenantId: job.TenantId, Occurrence: occurrence, Scheduled: slot != nil, Status: "queued", Command: job.Command, RetryTimes: job.RetryTimes, RetryInterval: job.RetryInterval, Timeout: job.Timeout, NextAttemptAt: now}
		if count >= limit {
			if slot == nil {
				return errors.New("job.busy")
			}
			run.Status = "skipped"
			run.FinishedAt = &now
		}
		if err := tx.Create(run).Error; err != nil {
			return err
		}
		if slot != nil {
			return tx.Model(&system.SysJob{}).Where("id = ?", job.Id).UpdateColumn("scheduled_at", *slot).Error
		}
		return nil
	})
	if err == nil && run != nil && run.Status == "queued" {
		js.wakeWorkers()
	}
	return run, err
}
func (js *JobScheduler) SubmitJob(job *system.SysJob) (*system.SysJobRun, error) {
	return js.submit(job, nil)
}
func (js *JobScheduler) ExecuteScheduledJob(job *system.SysJob, slot time.Time) error {
	run, err := js.submit(job, &slot)
	if err != nil || run == nil || run.Status == "skipped" {
		return err
	}
	return js.runAttempt(context.Background(), run.ID)
}
func (js *JobScheduler) ExecuteJob(job *system.SysJob) error {
	run, err := js.SubmitJob(job)
	if err != nil {
		return err
	}
	return js.runAttempt(context.Background(), run.ID)
}

func invokeJob(handler jobHandler, ctx context.Context, js *JobScheduler) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("handler panic: %v", recovered)
		}
	}()
	return handler(ctx, js)
}

// A savepoint rolls back a failed handler without losing the durable retry
// record. The occurrence lock fences all replicas for the whole attempt.
func (js *JobScheduler) runAttempt(parent context.Context, id string) error {
	// The outer transaction has time to persist a timeout and its next retry.
	ctx, cancel := context.WithTimeout(parent, 3610*time.Second)
	defer cancel()
	var attemptErr error
	err := platformMaintenanceDB(ctx).Transaction(func(tx *gorm.DB) error {
		var run system.SysJobRun
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).First(&run, "id = ?", id).Error; err != nil {
			return err
		}
		if run.Status != "queued" && run.Status != "retrying" || time.Now().Before(run.NextAttemptAt) {
			return nil
		}
		var job system.SysJob
		if err := tx.First(&job, run.JobID).Error; err != nil {
			attemptErr = err
		}
		if attemptErr == nil && job.TenantId != run.TenantId {
			attemptErr = errors.New("job ownership changed")
		}
		if attemptErr == nil && run.Scheduled && job.Status != 1 {
			attemptErr = errors.New("job stopped")
		}
		if attemptErr == nil && run.TenantId != global.PlatformTenantID {
			var tenant system.SysTenant
			if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&tenant, "tenant_id = ?", run.TenantId).Error; err != nil {
				attemptErr = err
			} else if !tenant.Status || util.IsTenantExpired(tenant.ExpireTime) {
				attemptErr = errors.New("tenantDisabled")
			}
		}
		run.Attempts++
		terminal := attemptErr != nil
		if attemptErr == nil {
			timeout := run.Timeout
			if timeout == 0 {
				timeout = 60
			}
			if err := tx.SavePoint("job_effects").Error; err != nil {
				return err
			}
			handlerCtx, stop := context.WithTimeout(jobContext(ctx, &job), time.Duration(timeout)*time.Second)
			// PostgreSQL statement_timeout preserves the connection and savepoint after a slow SQL statement.
			if err := tx.Exec("SELECT set_config('statement_timeout', ?, true)", strconv.FormatUint(uint64(timeout)*1000, 10)).Error; err != nil {
				stop()
				return err
			}
			handlerCtx = context.WithValue(handlerCtx, global.ContextDBKey, tx.WithContext(jobContext(ctx, &job)))
			handler, ok := jobRegistry[run.Command]
			if !ok {
				attemptErr = errors.New("unregistered job command")
				terminal = true
			} else {
				attemptErr = invokeJob(handler, handlerCtx, js)
			}
			if attemptErr == nil {
				attemptErr = handlerCtx.Err()
			}
			stop()
			if attemptErr != nil {
				if err := tx.RollbackTo("job_effects").Error; err != nil {
					return err
				}
			}
		}
		if err := tx.Exec("SELECT set_config('statement_timeout', '0', true)").Error; err != nil {
			return err
		}
		now := time.Now().UTC()
		updates := map[string]any{"attempts": run.Attempts}
		if attemptErr == nil {
			updates["status"] = "succeeded"
			updates["last_error"] = ""
			updates["finished_at"] = now
			if err := tx.Model(&system.SysJob{}).Where("id = ?", run.JobID).UpdateColumns(map[string]any{"times": gorm.Expr("times + 1"), "last_time": now}).Error; err != nil {
				return err
			}
		} else {
			runes := []rune(attemptErr.Error())
			if len(runes) > 256 {
				runes = runes[:256]
			}
			detail := string(runes)
			updates["last_error"] = detail
			if terminal || run.Attempts > run.RetryTimes {
				updates["status"] = "failed"
				updates["finished_at"] = now
			} else {
				updates["status"] = "retrying"
				interval := run.RetryInterval
				if interval == 0 {
					interval = 1
				}
				delay := time.Duration(interval) * time.Second * time.Duration(1<<(run.Attempts-1))
				if delay > 300*time.Second {
					delay = 300 * time.Second
				}
				updates["next_attempt_at"] = now.Add(delay)
			}
		}
		return tx.Model(&system.SysJobRun{}).Where("id = ?", id).UpdateColumns(updates).Error
	})
	if err != nil {
		return err
	}
	return attemptErr
}
func (js *JobScheduler) dispatchPending(stop <-chan struct{}) {
	var runs []system.SysJobRun
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := platformMaintenanceDB(ctx).Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).Where(activeRuns).Where("next_attempt_at <= ?", time.Now()).Order("next_attempt_at,id").Limit(32).Find(&runs).Error; err != nil {
		if global.Log != nil {
			global.Log.Warn("任务队列读取失败", zap.Error(err))
		}
		return
	}
	for _, run := range runs {
		select {
		case <-stop:
			return
		default:
		}
		select {
		case js.workers <- struct{}{}:
		default:
			return
		}
		js.workerWG.Add(1)
		go func(id string) {
			defer js.workerWG.Done()
			defer func() { <-js.workers; js.wakeWorkers() }()
			if err := js.runAttempt(context.Background(), id); err != nil && global.Log != nil && !strings.Contains(err.Error(), "55P03") {
				global.Log.Warn("任务执行失败", zap.String("runId", id), zap.Error(err))
			}
		}(run.ID)
	}
}

func (js *JobScheduler) wakeWorkers() {
	select {
	case js.workReady <- struct{}{}:
	default:
	}
}
