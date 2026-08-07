package config

// System 系统配置
type System struct {
	Port int `yaml:"port" mapstructure:"port"` // 服务端口
}
