package config

// RateLimit 租户级聚合请求限流配置。
// 鉴权后按 tenant_id 计数的固定窗口预算，用于单租户突发流量不至于拖垮其它租户；
// 平台域（超管/无租户上下文）用更高的独立档。enable=false 时整项不生效。
type RateLimit struct {
	Enable         bool `yaml:"enable" mapstructure:"enable"`
	WindowSeconds  int  `yaml:"window_seconds" mapstructure:"window_seconds"`
	BusinessPerSec int  `yaml:"business_per_second" mapstructure:"business_per_second"`
	PlatformPerSec int  `yaml:"platform_per_second" mapstructure:"platform_per_second"`
}
