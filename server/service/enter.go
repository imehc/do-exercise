package service

import "github.com/imehc/do-exercise/server/service/system"

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup system.ServiceGroup
}
