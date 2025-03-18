package main

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/initialize"
	"github.com/imehc/do-exercise/server/middleware"
	"github.com/labstack/echo/v4"
)

func init() {
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitGorm()
	initialize.InitRedis()
}

func main() {

	// 创建Echo实例
	e := echo.New()

	// 添加日志中间件
	e.Use(middleware.Logger)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", global.Config.System.Port)
	e.Logger.Fatal(e.Start(serverAddr))
}
