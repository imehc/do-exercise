package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/util"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig(configFile string) {
	// 创建viper实例
	v := viper.New()

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		util.Exit("配置文件不存在: ", err)
	}

	// 设置配置文件
	v.SetConfigFile(configFile)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnv(v, map[string]string{
		"oss.bucket_name":    "OSS_BUCKET_NAME",
		"oss.host":           "OSS_HOST",
		"oss.port":           "OSS_PORT",
		"oss.access_key":     "OSS_APP_ACCESS_KEY",
		"oss.secret_key":     "OSS_APP_SECRET_KEY",
		"oss.presigned_host": "OSS_PRESIGNED_HOST",
		"database.host":      "POSTGRES_HOST",
		"database.username":  "POSTGRES_USER",
		"database.password":  "POSTGRES_PASSWORD",
		"database.database":  "POSTGRES_DB",
		"database.port":      "POSTGRES_PORT",
		"redis.host":         "REDIS_HOST",
		"redis.port":         "REDIS_PORT",
		"redis.password":     "REDIS_PASSWORD",
		"redis.database":     "REDIS_DATABASE",
		"email.host":         "EMAIL_HOST",
		"email.port":         "EMAIL_PORT",
		"email.user":         "EMAIL_USER",
		"email.pass":         "EMAIL_PASS",
	})

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			util.Exit("配置文件未找到: ", err)
		} else {
			// 配置文件存在但产生了其他错误
			util.Exit("配置文件错误: ", err)
		}
	}

	// 将配置映射到全局变量
	if err := v.Unmarshal(&global.Config); err != nil {
		util.Exit("配置映射失败: ", err)
	}

	// 启动时一次性解析并校验 token 时长配置，避免请求期解析失败时生成 0 TTL 的永久 token。
	util.InitAuthDurations(global.Config.Auth.AccessExpireTime, global.Config.Auth.RefreshExpireTime)

	fmt.Println("配置初始化完成")
}

func bindEnv(v *viper.Viper, envs map[string]string) {
	for key, env := range envs {
		if err := v.BindEnv(key, env); err != nil {
			util.Exit("环境变量绑定失败: ", err)
		}
	}
}
