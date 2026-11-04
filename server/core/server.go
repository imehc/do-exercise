package core

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/middleware"
	"github.com/imehc/do-exercise/server/model"
	"github.com/imehc/do-exercise/server/model/system"
	"github.com/imehc/do-exercise/server/router"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	r := router.Run()

	// 启动操作日志批量落库 worker
	middleware.StartOperationLogWorker()

	// syncApi 由独立 goroutine 执行，用 WaitGroup 跟踪，优雅关闭时先等它收尾，
	// 避免进程退出截断尚未提交的 DB 写入。
	var syncApiWG sync.WaitGroup
	syncApiWG.Add(1)
	go func() {
		defer syncApiWG.Done()
		syncApi(r.Routes())
	}()

	// host := global.Config.System.Host
	addr := fmt.Sprintf(":%d", global.Config.System.Port)
	s := initServer(addr, r)

	// 监听失败（端口占用等）时以非零码退出；优雅关闭时 err 为 nil 或 ErrServerClosed，正常返回
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		global.Log.Fatal("server listen failed", zap.Error(err))
	}
	syncApiWG.Wait()
}

// syncApi 同步api到数据库
func syncApi(routes gin.RoutesInfo) {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// Route discovery runs on every replica. Lock before reading to avoid duplicate IDs and paths.
		if err := tx.Exec("LOCK TABLE sys_api IN SHARE ROW EXCLUSIVE MODE").Error; err != nil {
			return err
		}
		var existing []system.SysApi
		if err := tx.Find(&existing).Error; err != nil {
			return err
		}
		known := make(map[string]bool, len(existing))
		for _, api := range existing {
			known[api.Method+" "+api.Path] = true
		}
		var maxID uint
		if err := tx.Model(&system.SysApi{}).Select("COALESCE(MAX(id),0)").Scan(&maxID).Error; err != nil {
			return err
		}
		for _, route := range routes {
			if known[route.Method+" "+route.Path] {
				continue
			}
			maxID++
			api := system.SysApi{IdWrapper: model.IdWrapper{Id: maxID}, Path: route.Path, Method: route.Method}
			if err := tx.Create(&api).Error; err != nil {
				return err
			}
		}
		return tx.Exec("SELECT setval(pg_get_serial_sequence('sys_api','id'), GREATEST((SELECT COALESCE(MAX(id),1) FROM sys_api),(SELECT last_value FROM sys_api_id_seq)),true)").Error
	})
	if err != nil {
		global.Log.Error("同步api失败", zap.Error(err))
	}
}
