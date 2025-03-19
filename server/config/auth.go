package config

type Auth struct {
	Domain            string `yaml:"domain" mapstructure:"domain"`                           // 域名
	AccessSecret      string `yaml:"access_secret" mapstructure:"access_secret"`             // 访问密钥
	AccessExpireTime  string `yaml:"access_expire_time" mapstructure:"access_expire_time"`   // 访问过期时间
	RefreshSecret     string `yaml:"refresh_secret" mapstructure:"refresh_secret"`           // 刷新密钥
	RefreshExpireTime string `yaml:"refresh_expire_time" mapstructure:"refresh_expire_time"` // 刷新过期时间
}
