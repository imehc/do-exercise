package config

type Oss struct {
	Host          string `yaml:"host" mapstructure:"host"`                     // oss服务器地址
	Port          int    `yaml:"port" mapstructure:"port"`                     // oss服务器端口
	AccessKey     string `yaml:"access_key" mapstructure:"access_key"`         // oss服务器用户名
	SecretKey     string `yaml:"secret_key" mapstructure:"secret_key"`         // oss服务器密码
	BucketName    string `yaml:"bucket_name" mapstructure:"bucket_name"`       // oss服务器桶名
	Expires       int    `yaml:"expires" mapstructure:"expires"`               // 预签名有效期，单位秒
	PresignedHost string `yaml:"presigned_host" mapstructure:"presigned_host"` // 预签名地址
	Secure        bool   `yaml:"secure" mapstructure:"secure"`                 // 是否使用HTTPS连接
}
