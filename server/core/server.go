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

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	s := initServer(address, r)
	fmt.Printf(`swagger文档地址:http://127.0.0.1%s/swagger/index.html`, address)

	global.LOG.Error(s.ListenAndServe().Error())
}
