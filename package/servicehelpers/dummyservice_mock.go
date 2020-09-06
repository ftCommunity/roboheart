package servicehelpers

import "github.com/ftCommunity-roboheart/roboheart/package/service"

type dummyService struct {
	deps, adeps []string
}

func (dummyService) Init(map[string]service.Service, service.LoggerFunc, service.ErrorFunc) {}
func (dummyService) Stop()                                                                  {}
func (dummyService) Name() string {
	return ""
}

func (d dummyService) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{
		Deps:  d.deps,
		ADeps: d.adeps,
	}
}
