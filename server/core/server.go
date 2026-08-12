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
	"github.com/samber/lo"
	"go.uber.org/zap"
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
	// 1. 获取所有唯一标识
	var uniqueFields []string
	for _, route := range routes {
		uniqueFields = append(uniqueFields, route.Path)
	}
	// 2. 查询已存在的记录
	var existingApis []system.SysApi
	if err := global.DB.Where("path IN ?", uniqueFields).Find(&existingApis).Error; err != nil {
		global.Log.Error("查询已存在的api失败", zap.Error(err))
	}
	// 3. 构建已存在记录的map
	existingMap := make(map[string]system.SysApi)
	for _, api := range existingApis {
		existingMap[api.Path] = api
	}
	// 4. 筛选出需要插入的记录
	var insertApis []system.SysApi
	for _, route := range routes {
		if _, ok := existingMap[route.Path]; !ok {
			insertApis = append(insertApis, system.SysApi{
				Path:   route.Path,
				Method: route.Method,
			})
		}
	}
	// 5. 插入新记录
	if len(insertApis) > 0 {
		db := global.DB
		var maxID int64
		db.Table("sys_api").Select("COALESCE(MAX(id), 0)").Row().Scan(&maxID)
		apis := lo.Map(insertApis, func(item system.SysApi, index int) *system.SysApi {
			return &system.SysApi{
				IdWrapper: model.IdWrapper{
					Id: uint(maxID) + uint(index) + 1,
				},
				Path:   item.Path,
				Method: item.Method,
			}
		})

		// 失败必须回滚。原先是 defer tx.Commit()，即使 Create 报错也照样提交，
		// 会把部分写入的路由记录持久化下来。
		tx := db.Begin()
		if err := tx.Create(apis).Error; err != nil {
			tx.Rollback()
			global.Log.Error("同步api失败", zap.Error(err))
			return
		}
		if err := tx.Commit().Error; err != nil {
			global.Log.Error("同步api提交失败", zap.Error(err))
		}
	}
}
