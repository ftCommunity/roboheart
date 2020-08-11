package servicemanager

import (
	"github.com/ftCommunity/roboheart/internal/service"
)

type ServiceState struct {
	sm      *ServiceManager
	name    string
	service struct {
		base          service.Service
		emerstoppable service.EmergencyStoppableService
		depending     service.DependingService
		adddepending  service.AddDependingService
	}
	running bool
	deps    struct {
		service.ServiceDependencies
		deps   ssmap
		rdeps  map[string]*ServiceState
		adeps  map[string]adep
		radeps map[string]*adep
	}
	logger service.LoggerFunc
	error  service.ErrorFunc
}

type ssmap map[string]*ServiceState

func (ssm *ssmap) toServiceMap() map[string]service.Service {
	svcs := make(map[string]service.Service)
	for sn, s := range *ssm {
		svcs[sn] = s.service.base
	}
	return svcs
}

type adep struct {
	*ServiceState
	set bool
}

func (ss *ServiceState) getBase() service.Service {
	return ss.service.base
}

func (ss *ServiceState) checkDepsReady() bool {
	for _, ds := range ss.deps.deps {
		if !ds.running {
			return false
		}
	}
	return true
}

func (ss *ServiceState) tryRun() bool {
	if ss.running {
		return true
	}
	if !ss.checkDepsReady() {
		return false
	}
	ss.logger("Starting...")
	ss.getBase().Init(ss.deps.deps.toServiceMap(), ss.logger, ss.error)
	ss.logger("Started successfully")
	ss.running = true
	return true
}

func (ss *ServiceState) setReadyAdeps() bool {
	if !ss.running {
		panic("cannot set additional dependencies on stopped service")
	}
	if ss.service.adddepending == nil {
		return true
	}
	done := true
	adeps := make(map[string]service.Service)
	for _, ds := range ss.deps.adeps {
		if !ds.set {
			if ds.running {
				adeps[ds.name] = ds.getBase()
				ds.set = true
				rd := ds.deps.radeps[ss.name]
				rd.set = true
			} else {
				done = false
			}
		}
	}
	if len(adeps) > 0 {
		ss.service.adddepending.SetAdditionalDependencies(adeps)
	}
	return done
}

func (ss *ServiceState) tryStop() bool {
	if !ss.running {
		return true
	}
	for _, rds := range ss.deps.radeps {
		if rds.set {
			rds.service.adddepending.UnsetAdditionalDependencies([]string{ss.name})
			ts := rds.deps.adeps[ss.name]
			ts.set = false
			rds.set = false
		}
	}
	for _, ds := range ss.deps.rdeps {
		if ds.running {
			return false
		}
	}
	ss.logger("Stopping...")
	ss.getBase().Stop()
	ss.logger("Stopped successfully")
	ss.running = false
	return true
}

func (ss *ServiceState) emerstop() {
	if emer := ss.service.emerstoppable; emer != nil {
		emer.EmergencyStop()
	}
}

func (ss *ServiceState) loadDepData() {
	for _, sn := range ss.deps.Deps {
		ss.deps.deps[sn] = ss.sm.services[sn]
		ss.sm.services[sn].deps.rdeps[ss.name] = ss
	}
	for _, sn := range ss.deps.ADeps {
		ss.deps.adeps[sn] = adep{
			ServiceState: ss.sm.services[sn],
		}
		ss.sm.services[sn].deps.radeps[ss.name] = &adep{
			ServiceState: ss,
		}
	}
}

func newServiceState(sm *ServiceManager, s service.Service) *ServiceState {
	ss := &ServiceState{}
	ss.name = s.Name()
	ss.service.base = s
	if es, ok := s.(service.EmergencyStoppableService); ok {
		ss.service.emerstoppable = es
	}
	if ds, ok := s.(service.DependingService); ok {
		ss.service.depending = ds
		ss.deps.ServiceDependencies = ds.Dependencies()
		if ads, ok := s.(service.AddDependingService); ok {
			ss.service.adddepending = ads
		}
	}
	ss.logger = sm.genServiceLogger(ss.name)
	ss.error = sm.genServiceError(ss.name)
	ss.sm = sm
	ss.deps.deps = make(ssmap)
	ss.deps.rdeps = make(map[string]*ServiceState)
	ss.deps.adeps = make(map[string]adep)
	ss.deps.radeps = make(map[string]*adep)
	return ss
}
