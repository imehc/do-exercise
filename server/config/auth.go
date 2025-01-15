package config

type AuthConfig struct {
	ClientID               string `mapstructure:"client_id" json:"client_id" yaml:"client_id"`
	ClientSecret           string `mapstructure:"client_secret" json:"client_secret" yaml:"client_secret"`
	Domain                 string `mapstructure:"domain" json:"domain" yaml:"domain"`
	AccessExpire           string `mapstructure:"access_expire_time" json:"access_expire_time" yaml:"access_expire_time"`
	RefreshExpire          string `mapstructure:"refresh_expire_time" json:"refresh_expire_time" yaml:"refresh_expire_time"`
	IsGenerateRefreshToken bool   `mapstructure:"is_generate_refresh_token" json:"is_generate_refresh_token" yaml:"is_generate_refresh_token"`
}
