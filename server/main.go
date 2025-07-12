package main

import (
	_ "net/http/pprof"

	"github.com/imehc/do-exercise/server/cmd"
)

func main() {
	cmd.Execute()
}
