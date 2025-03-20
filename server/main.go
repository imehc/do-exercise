package main

import (
	"github.com/imehc/do-exercise/server/core"
	"github.com/imehc/do-exercise/server/internal"
)

func init() {
	internal.InitConfig()
	internal.InitLogger()
	internal.InitGorm()
	internal.InitRedis()
	internal.InitI18n()
}

func main() {
	core.RunServer()
}
