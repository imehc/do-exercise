package router

import (
	"github.com/imehc/do-exercise/server/router/common"
	"github.com/imehc/do-exercise/server/router/normal"
	"github.com/imehc/do-exercise/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System system.RouterGroup
	Common common.RouterGroup
	Normal normal.RouterGroup
}
