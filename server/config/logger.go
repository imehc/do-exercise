package config

type Logger struct {
	Level      string `yaml:"level" mapstructure:"level"`           // 日志级别
	Directory  string `yaml:"directory" mapstructure:"directory"`   // 日志目录
	MaxSize    int    `yaml:"maxSize" mapstructure:"maxSize"`       // 单个文件最大大小，单位KB
	MaxBackups int    `yaml:"maxBackups" mapstructure:"maxBackups"` // 最大保留文件数
	MaxAge     int    `yaml:"maxAge" mapstructure:"maxAge"`         // 最大保留天数
	Compress   bool   `yaml:"compress" mapstructure:"compress"`     // 是否压缩
}