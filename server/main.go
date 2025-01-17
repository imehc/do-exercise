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

// docs https://goswagger.io/go-swagger/
// @title                       					API接口文档
// @version                     					v0.1.0
// @description                 					API接口文档
// @securitydefinitions.bearerauth        bearer
func main() {
	core.RunServer()
}
