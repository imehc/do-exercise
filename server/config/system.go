package config

// System 系统配置
type System struct {
	Port int `yaml:"port" mapstructure:"port"` // 服务端口
	// TrustedProxies 受信任的反向代理 CIDR 列表。
	//
	// gin 默认信任所有代理（0.0.0.0/0、::/0），此时 c.ClientIP() 会直接采信
	// 客户端自带的 X-Forwarded-For，导致限流、登录锁定被绕过，审计日志 IP 可伪造。
	// 留空表示不信任任何代理，ClientIP() 退回使用 TCP 连接的真实对端地址。
	// 部署在 nginx/LB 之后时，填其所在网段，例如 ["166.6.0.0/16"]。
	TrustedProxies []string `yaml:"trusted_proxies" mapstructure:"trusted_proxies"`
}
