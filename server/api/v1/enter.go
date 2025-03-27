package v1

import (
	"github.com/imehc/do-exercise/server/api/v1/common"
	"github.com/imehc/do-exercise/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup
	CommonApiGroup common.ApiGroup
}
