package config

// Redis Redis配置
type Redis struct {
	Host     string       `yaml:"host" mapstructure:"host"`         // Redis主机地址
	Port     int          `yaml:"port" mapstructure:"port"`         // Redis端口
	Password string       `yaml:"password" mapstructure:"password"` // Redis密码
	Database int          `yaml:"database" mapstructure:"database"` // Redis数据库
	Pool     RedisPool    `yaml:"pool" mapstructure:"pool"`         // 连接池配置
	Timeout  RedisTimeout `yaml:"timeout" mapstructure:"timeout"`   // 超时配置
}

// RedisPool Redis连接池配置
type RedisPool struct {
	MaxConnections     int `yaml:"max_connections" mapstructure:"max_connections"`           // 最大连接数
	MinIdleConnections int `yaml:"min_idle_connections" mapstructure:"min_idle_connections"` // 最小空闲连接数
	MaxIdleTime        int `yaml:"max_idle_time" mapstructure:"max_idle_time"`               // 最大空闲时间(秒)
}

// RedisTimeout Redis超时配置
type RedisTimeout struct {
	Connect int `yaml:"connect" mapstructure:"connect"` // 连接超时时间(秒)
	Read    int `yaml:"read" mapstructure:"read"`       // 读取超时时间(秒)
	Write   int `yaml:"write" mapstructure:"write"`     // 写入超时时间(秒)
}
