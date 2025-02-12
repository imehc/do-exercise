package internal

import (
	"fmt"

	"github.com/imehc/do-exercise/server/config"
	"github.com/imehc/do-exercise/server/global"
	"gorm.io/gorm/logger"
)

type Writer struct {
	config config.Datebase
	writer logger.Writer
}

func NewWriter(config config.Datebase) *Writer {
	return &Writer{config: config}
}

func (c *Writer) Printf(message string, data ...any) {
	// 当有日志时候均需要输出到控制台
	fmt.Printf(message, data...)

	// 当开启了zap的情况，会打印到日志记录
	if c.config.LogZap {
		switch c.config.LogLevel() {
		case logger.Silent:
			global.LOG.Debug(fmt.Sprintf(message, data...))
		case logger.Error:
			global.LOG.Error(fmt.Sprintf(message, data...))
		case logger.Warn:
			global.LOG.Warn(fmt.Sprintf(message, data...))
		case logger.Info:
			global.LOG.Info(fmt.Sprintf(message, data...))
		default:
			global.LOG.Info(fmt.Sprintf(message, data...))
		}
		return
	}
}
