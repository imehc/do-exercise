package config

type Minio struct {
	Host       string `yaml:"host" mapstructure:"host"`               // minio服务器地址
	Port       string `yaml:"port" mapstructure:"port"`               // minio服务器端口
	AccessKey  string `yaml:"access_key" mapstructure:"access_key"`   // minio服务器用户名
	SecretKey  string `yaml:"secret_key" mapstructure:"secret_key"`   // minio服务器密码
	BucketName string `yaml:"bucket_name" mapstructure:"bucket_name"` // minio服务器桶名
	Expires    int    `yaml:"expires" mapstructure:"expires"`         // 预签名有效期，单位秒
}
