package config

type Auth struct {
	AccessExpireTime  string `yaml:"access_expire_time" mapstructure:"access_expire_time"`   // 访问过期时间
	RefreshExpireTime string `yaml:"refresh_expire_time" mapstructure:"refresh_expire_time"` // 刷新过期时间
}
