package core

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/router"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {

	r := router.Run()

	// host := global.Config.System.Host
	addr := fmt.Sprintf(":%d", global.Config.System.Port)
	s := initServer(addr, r)

	global.Log.Error(s.ListenAndServe().Error())
}
