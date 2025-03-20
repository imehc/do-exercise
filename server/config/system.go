package config

// System 系统配置
type System struct {
	Port int    `yaml:"port" mapstructure:"port"` // 服务端口
	Host string `yaml:"host" mapstructure:"host"` // 服务地址
}
