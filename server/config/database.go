package config

// Database PostgreSQL数据库配置
type Database struct {
	Host     string `yaml:"host" mapstructure:"host"`         // 数据库主机地址
	Port     int    `yaml:"port" mapstructure:"port"`         // 数据库端口
	Database string `yaml:"database" mapstructure:"database"` // 数据库名称
	Username string `yaml:"username" mapstructure:"username"` // 数据库用户名
	Password string `yaml:"password" mapstructure:"password"` // 数据库密码
	Pool     Pool   `yaml:"pool" mapstructure:"pool"`         // 连接池配置
}

// Pool 数据库连接池配置
type Pool struct {
	MaxConnections    int `yaml:"max_connections" mapstructure:"max_connections"`       // 最大连接数
	MinConnections    int `yaml:"min_connections" mapstructure:"min_connections"`       // 最小连接数
	MaxIdleTime       int `yaml:"max_idle_time" mapstructure:"max_idle_time"`           // 最大空闲时间(秒)
	ConnectionTimeout int `yaml:"connection_timeout" mapstructure:"connection_timeout"` // 连接超时时间(秒)
}
