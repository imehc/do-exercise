package core

import (
	"fmt"

	"github.com/imehc/do-exercise/server/global"
	"github.com/imehc/do-exercise/server/initialize"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {

	r := initialize.InitRouter()

	host := global.CONFIG.System.Host
	addr := fmt.Sprintf(":%d", global.CONFIG.System.Port)
	s := initServer(addr, r)
	fmt.Printf(`swagger文档地址:http://%s%s/swagger/index.html`, host, addr)

	global.LOG.Error(s.ListenAndServe().Error())
}
