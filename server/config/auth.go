package config

type Auth struct {
	AccessExpireTime  string `yaml:"access_expire_time" mapstructure:"access_expire_time"`   // 访问过期时间
	RefreshExpireTime string `yaml:"refresh_expire_time" mapstructure:"refresh_expire_time"` // 刷新过期时间
	LoginMaxAttempts  int    `yaml:"login_max_attempts" mapstructure:"login_max_attempts"`   // 登录失败锁定阈值（0 使用默认 5）
	LoginLockMinutes  int    `yaml:"login_lock_minutes" mapstructure:"login_lock_minutes"`   // 登录失败锁定时长（分钟，0 使用默认 5）
}
