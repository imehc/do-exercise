package main

import (
	"github.com/imehc/do-exercise/server/core"
	"github.com/imehc/do-exercise/server/global/shared"
	"github.com/imehc/do-exercise/server/internal"
)

func init() {
	internal.InitConfig()
	internal.InitLogger()
	internal.InitGorm()
	internal.InitRedis()
	internal.InitI18n()
	internal.InitCasbin()
}

func main() {
	defer shared.RSACrypto.Stop()

	core.RunServer()
}
