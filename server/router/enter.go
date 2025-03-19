package router

import (
	"github.com/imehc/do-exercise/server/router/normal"
	"github.com/imehc/do-exercise/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System system.RouterGroup
	Normal normal.RouterGroup
}
