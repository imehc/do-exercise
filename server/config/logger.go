package config

type Logger struct {
	Level      string `yaml:"level" mapstructure:"level"`             // 日志级别
	Directory  string `yaml:"directory" mapstructure:"directory"`     // 日志目录
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`       // 单个文件最大大小，单位KB
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"` // 最大保留文件数
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`         // 最大保留天数
	Compress   bool   `yaml:"compress" mapstructure:"compress"`       // 是否压缩
	BucketName string `yaml:"bucket_name" mapstructure:"bucket_name"` // 桶名称
	Expires    int    `yaml:"expires" mapstructure:"expires"`         // 过期时间，单位秒
}
