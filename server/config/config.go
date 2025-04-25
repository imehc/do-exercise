package config

type Config struct {
	System   System   `yaml:"system" mapstructure:"system"`     // 系统配置
	Database Database `yaml:"database" mapstructure:"database"` // 数据库配置
	Redis    Redis    `yaml:"redis" mapstructure:"redis"`       // Redis配置
	Logger   Logger   `yaml:"logger" mapstructure:"logger"`     // 日志配置
	Auth     Auth     `yaml:"auth" mapstructure:"auth"`         // 认证配置
	Captcha  Captcha  `yaml:"captcha" mapstructure:"captcha"`   // 验证码配置
	Email    Email    `yaml:"email" mapstructure:"email"`       // 邮箱配置
	Minio    Minio    `yaml:"minio" mapstructure:"minio"`       // Minio配置
}
