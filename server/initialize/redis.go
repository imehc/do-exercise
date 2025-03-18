package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/imehc/do-exercise/server/global"
	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化Redis连接
func InitRedis() {
	// 创建Redis配置
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.Config.Redis.Host, global.Config.Redis.Port),
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.Database,

		// 连接池配置
		PoolSize:     global.Config.Redis.Pool.MaxConnections,
		MinIdleConns: global.Config.Redis.Pool.MinIdleConnections,
		MaxIdleConns: global.Config.Redis.Pool.MaxConnections,

		// 超时配置
		DialTimeout:  time.Duration(global.Config.Redis.Timeout.Connect) * time.Second,
		ReadTimeout:  time.Duration(global.Config.Redis.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(global.Config.Redis.Timeout.Write) * time.Second,
	}

	// 创建Redis客户端
	client := redis.NewClient(options)

	// 测试连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("连接Redis失败: %v", err))
	}

	// 将Redis客户端设置为全局变量
	global.Redis = client

	fmt.Println("Redis连接初始化成功")
}
