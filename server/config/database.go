package config

import (
	"strings"

	"gorm.io/gorm/logger"
)

type Datebase struct {
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               // 数据库端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // 高级配置
	DbName       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 数据库账号
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 数据库密码
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                               // 数据库地址
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                   // 是否开启Gorm全局日志
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`                      // 是否通过zap写入日志文件
}

func (d Datebase) LogLevel() logger.LogLevel {
	switch strings.ToLower(d.LogMode) {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}

// 基于配置文件获取 dsn
func (d *Datebase) Dsn() string {
	return "host=" + d.Path + " user=" + d.Username + " password=" + d.Password + " dbname=" + d.DbName + " port=" + d.Port + " " + d.Config
}

// 根据 dbname 生成 dsn
func (d *Datebase) LinkDsn(dbname string) string {
	return "host=" + d.Path + " user=" + d.Username + " password=" + d.Password + " dbname=" + dbname + " port=" + d.Port + " " + d.Config
}
