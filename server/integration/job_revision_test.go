//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/model/system"
	"gorm.io/gorm"
)

func TestJobScheduleRevisionAfterTimestampRoundTrip(t *testing.T) {
	u := principal(t)
	// Keep nanoseconds in the in-memory snapshot even on clocks with lower resolution.
	stamp := time.Date(2026, 11, 4, 10, 0, 0, 123456789, time.FixedZone("UTC+8", 8*60*60))
	job := &system.SysJob{
		Name:           "timestamp-round-trip",
		JobGroup:       "revision",
		TenantId:       u.TenantId,
		CronExpression: "0 0 0 * * *",
		Command:        "clean_empty_username_operation_logs",
		Status:         1,
	}
	must(t, bypass().Session(&gorm.Session{NowFunc: func() time.Time { return stamp }}).Create(job).Error)
	var stored system.SysJob
	must(t, bypass().First(&stored, job.Id).Error)
	if job.UpdatedAt.Nanosecond() != stamp.Nanosecond() || stored.UpdatedAt.UnixNano()%1000 != 0 || stored.UpdatedAt.Equal(job.UpdatedAt) {
		t.Fatal("fixture must retain nanoseconds in memory and lose them in PostgreSQL")
	}

	js := &shared.JobScheduler{}
	slot := time.Now().UTC().Truncate(time.Second)
	must(t, js.ExecuteScheduledJob(job, slot))
	must(t, js.ExecuteScheduledJob(job, slot))
	must(t, bypass().First(&stored, job.Id).Error)
	if stored.Times != 1 {
		t.Fatalf("same occurrence ran %d times after the timestamp round trip", stored.Times)
	}

	// A real update at PostgreSQL's precision must still invalidate the snapshot.
	must(t, bypass().Model(&system.SysJob{}).Where("id = ?", job.Id).UpdateColumn("updated_at", stored.UpdatedAt.Add(time.Microsecond)).Error)
	if err := js.ExecuteScheduledJob(job, slot.Add(time.Second)); err == nil {
		t.Fatal("outdated timestamp snapshot accepted")
	}
	must(t, bypass().First(&stored, job.Id).Error)
	must(t, js.ExecuteScheduledJob(&stored, slot.Add(time.Second)))
	must(t, bypass().First(&stored, job.Id).Error)
	if stored.Times != 2 {
		t.Fatalf("updated schedule ran %d times", stored.Times)
	}

	// Cron changes must be detected even when updated_at is unchanged.
	must(t, bypass().Model(&system.SysJob{}).Where("id = ?", job.Id).UpdateColumn("cron_expression", "0 30 0 * * *").Error)
	if err := js.ExecuteScheduledJob(&stored, slot.Add(2*time.Second)); err == nil {
		t.Fatal("outdated cron snapshot accepted")
	}
	var count int64
	must(t, bypass().Model(&system.SysJobRun{}).Where("job_id = ?", job.Id).Count(&count).Error)
	if count != 2 {
		t.Fatalf("rejected or duplicate snapshots created %d occurrences", count)
	}
}
