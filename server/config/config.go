package config

type Config struct {
	Auth    AuthConfig `mapstructure:"auth" json:"auth" yaml:"auth"`
	Zap     Zap        `mapstructure:"zap" json:"zap" yaml:"zap"`
	System  System     `mapstructure:"system" json:"system" yaml:"system"`
	Redis   Redis      `mapstructure:"redis" json:"redis" yaml:"redis"`
	Captcha Captcha    `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}
