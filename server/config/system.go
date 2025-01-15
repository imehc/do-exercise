package config

type System struct {
	RouterPrefix string `mapstructure:"router_prefix" json:"router_prefix" yaml:"router_prefix"`
	Addr         int    `mapstructure:"addr" json:"addr" yaml:"addr"` // 端口值
	UseRedis     bool   `mapstructure:"use_redis" json:"use_redis" yaml:"use_redis"`
}
