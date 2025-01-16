package config

type System struct {
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	RouterPrefix string `mapstructure:"router_prefix" json:"router_prefix" yaml:"router_prefix"`
	UseRedis     bool   `mapstructure:"use_redis" json:"use_redis" yaml:"use_redis"`
}
