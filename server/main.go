package main

import (
	"github.com/imehc/do-exercise/server/internal"
	"github.com/imehc/do-exercise/server/router"
)

func init() {
	internal.InitConfig()
	internal.InitLogger()
	// internal.InitGorm()
	// internal.InitRedis()
	internal.InitI18n()
}

func main() {
	router.RunServer()
}
