package main

import (
	"flag"
	"fmt"

	"github.com/imehc/do-exercise/server/core"
	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/initialize"
	_ "go.uber.org/automaxprocs"
)

// 定义命令行参数
var (
	initData = flag.Bool("init-data", false, "初始化数据库数据")
)

func init() {
	// 解析命令行参数
	flag.Parse()

	initialize.InitViper()
	initialize.InitZap()
	initialize.InitCache()
	initialize.InitPostgres()
	initialize.InitCasbin()
	initialize.InitAuth()
	initialize.InitOther()
}

// docs https://goswagger.io/go-swagger/
// @title                       					API接口文档
// @version                     					v0.1.0
// @description                 					API接口文档
// @securitydefinitions.bearerauth        bearer
func main() {

	if global.DB != nil {
		initialize.InitRegisterTable()
		fmt.Println("数据库初始化成功")

		// 如果指定了-init-data参数，执行数据初始化
		if *initData {
			initialize.InitData()
		}

		db, _ := global.DB.DB()
		defer db.Close()
	}

	core.RunServer()
}
