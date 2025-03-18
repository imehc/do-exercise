package config

type Config struct {
	System   System   `yaml:"system" mapstructure:"system"`     // 系统配置
	Database Database `yaml:"database" mapstructure:"database"` // 数据库配置
	Redis    Redis    `yaml:"redis" mapstructure:"redis"`       // Redis配置
	Logger   Logger   `yaml:"logger" mapstructure:"logger"`     // 日志配置
}
