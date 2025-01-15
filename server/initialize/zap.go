package initialize

import (
	"fmt"
	"os"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/initialize/internal"
	"github.com/imehc/do-exercise/server/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitZap 初始化 zap 日志库。
// 该函数检查日志目录是否存在，如果不存在则创建该目录，
// 然后根据全局配置设置日志级别和是否显示调用行号等配置项。
func InitZap() error {
	// 检查日志目录是否存在，如果不存在则创建
	if ok, _ := utils.HasPathExists(global.CONFIG.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", global.CONFIG.Zap.Director)
		_ = os.Mkdir(global.CONFIG.Zap.Director, os.ModePerm)
	}

	// 从配置中获取所有日志级别
	levels := global.CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)

	// 遍历每个日志级别并创建相应的 zap core
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}

	// 创建一个新的 logger，使用配置的多个 core
	logger := zap.New(zapcore.NewTee(cores...))

	// 如果配置中启用了显示调用行号，则启用该选项
	if global.CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	// 设置全局 logger 并替换默认的全局 logger
	global.LOG = logger
	zap.ReplaceGlobals(global.LOG)

	return nil
}
