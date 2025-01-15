package main

import (
	"github.com/imehc/do-exercise/server/core"
	"github.com/imehc/do-exercise/server/initialize"
	_ "go.uber.org/automaxprocs"
)

func init() {
	initialize.InitViper()
	initialize.InitZap()
	initialize.InitCache()
	initialize.InitAuth()
	initialize.InitOther()
}

// @title                       API接口文档
// @version                     v0.1.0
// @description                 接口文档
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 令牌授权方式，格式：Bearer <token>
func main() {
	core.RunServer()
}
