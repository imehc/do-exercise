package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化日志
func InitLogger() {
	// 创建日志目录
	logConfig := global.Config.Logger
	if err := os.MkdirAll(logConfig.Directory, 0o755); err != nil {
		util.Exit("创建日志目录失败: ", err)
	}

	// 设置日志级别
	var level zapcore.Level
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建不同级别的核心
	cores := []zapcore.Core{}

	// 定义日志级别和对应的文件名
	logLevels := map[string]zapcore.Level{
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"debug": zapcore.DebugLevel,
	}

	// 为每个级别创建一个日志写入器
	for levelName, zapLevel := range logLevels {
		if level <= zapLevel {
			// 日志文件名格式：级别_年月日.log
			fileName := filepath.Join(logConfig.Directory, fmt.Sprintf("%s_%s.log",
				levelName, time.Now().Format("20060102")))

			// 配置日志切割。lumberjack 的 MaxSize 单位是 MB（此前按 KB 配置，
			// 10 万 KB 的写法让轮转阈值变成 10GB/文件，磁盘会被写满）
			lumberJackLogger := &lumberjack.Logger{
				Filename:   fileName,
				MaxSize:    logConfig.MaxSize,
				MaxBackups: logConfig.MaxBackups,
				MaxAge:     logConfig.MaxAge,
				Compress:   logConfig.Compress,
			}

			// 编码器配置
			encoderConfig := zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}

			// 创建核心
			writeSyncer := zapcore.AddSync(lumberJackLogger)
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				writeSyncer,
				zapLevel,
			)
			cores = append(cores, core)
		}
	}

	// 创建日志记录器
	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	// 替换全局日志记录器
	zap.ReplaceGlobals(logger)

	global.Log = logger

	fmt.Println("日志初始化完成")
}
