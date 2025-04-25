package config

type Email struct {
	Host string `yaml:"host" mapstructure:"host"` // 邮箱服务器地址
	Port int    `yaml:"port" mapstructure:"port"` // 邮箱服务器端口
	User string `yaml:"user" mapstructure:"user"` // 邮箱用户名
	Pass string `yaml:"pass" mapstructure:"pass"` // 邮箱密码
}
