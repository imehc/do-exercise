package config

type Captcha struct {
	ImgWidth       int `yaml:"img_width" mapstructure:"img_width"`              // 验证码宽度
	ImgHeight      int `yaml:"img_height" mapstructure:"img_height"`            // 验证码高度
	CaptchaTimeOut int `yaml:"captcha_timeout" mapstructure:"captcha_timeout"` // 验证码超时时间，单位：s(秒)
}
