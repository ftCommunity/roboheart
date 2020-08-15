package serviceadmin

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

type serviceadmin struct {
	acm     acm.ACM
	web     web.Web
	logger  service.LoggerFunc
	error   service.ErrorFunc
	manager service.ServiceManager
}

func (s *serviceadmin) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) {
	s.logger = logger
	s.error = e
	if err := servicehelpers.CheckMainDependencies(s, services); err != nil {
		e(err)
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{"acm": s.initSvcAcm}); err != nil {
		e(err)
	}
}

func (s *serviceadmin) Stop()        {}
func (s *serviceadmin) Name() string { return "serviceadmin" }
func (s *serviceadmin) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{Deps: []string{"acm"}, ADeps: []string{"web"}}
}
func (s *serviceadmin) SetAdditionalDependencies(services map[string]service.Service) {
	servicehelpers.InitializeAdditionalDependencies(services, servicehelpers.AdditionalServiceInitializers{"web": s.initSvcWeb})
}
func (s *serviceadmin) UnsetAdditionalDependencies(services []string) {
	servicehelpers.DeinitAdditionalDependencies(services, servicehelpers.AdditionalServiceDeinitializers{"web": s.deinitSvcWeb})
}

func (s *serviceadmin) SetServiceManager(manager service.ServiceManager) {
	s.manager = manager
}
