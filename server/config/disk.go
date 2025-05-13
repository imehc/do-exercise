package config

type Disk struct {
	MountPoint string `yaml:"mount_point" mapstructure:"mount_point"` //  挂载点
}
