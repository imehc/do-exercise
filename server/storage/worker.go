package storage

import (
	"context"
	"sync"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StartCollector returns a stop function that waits for the in-flight batch.
func StartCollector(db *gorm.DB) func() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		interval := global.Config.Oss.GCIntervalSeconds
		if interval <= 0 {
			interval = 60
		}
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for {
			batch, stop := context.WithTimeout(ctx, 30*time.Second)
			if err := ReconcileLegacy(batch, db, time.Now().UTC(), 100); err != nil && global.Log != nil {
				global.Log.Warn("对象元数据同步失败", zap.Error(err))
			}
			if _, err := Collect(batch, db, time.Now().UTC(), 100); err != nil && global.Log != nil {
				global.Log.Warn("对象回收失败", zap.Error(err))
			}
			maintenance := util.TenantBypass(db.WithContext(batch))
			cutoff := time.Now().Add(-30 * 24 * time.Hour)
			if err := maintenance.Where("finished_at < ?", cutoff).Delete(&system.SysJobRun{}).Error; err != nil && global.Log != nil {
				global.Log.Warn("任务记录清理失败", zap.Error(err))
			}
			if err := maintenance.Where("state = 'deleted' AND deleted_at < ?", cutoff).Delete(&system.SysFile{}).Error; err != nil && global.Log != nil {
				global.Log.Warn("文件记录清理失败", zap.Error(err))
			}
			stop()
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
	return func() { cancel(); wg.Wait() }
}
