package config

type Captcha struct {
	KeyLong        int `mapstructure:"key_long" json:"key_long" yaml:"key_long"`                      // 验证码长度
	ImgWidth       int `mapstructure:"img_width" json:"img_width" yaml:"img_width"`                   // 验证码宽度
	ImgHeight      int `mapstructure:"img_height" json:"img_height" yaml:"img_height"`                // 验证码高度
	CaptchaTimeOut int `mapstructure:"captcha_timeout" json:"captcha_timeout" yaml:"captcha_timeout"` // 验证码超时时间，单位：s(秒)
}
