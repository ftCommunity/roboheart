package servicemanager

import (
	"errors"
	"log"
	"sync"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/thoas/go-funk"
)

var (
	services = []service.Service{
	}
	coreservices = []service.Service{
		config.Service,
		acm.Service,
	}
)

type ServiceState struct {
	service service.Service
	core    bool
	running bool
	addset  bool
}

type ServiceManager struct {
	services map[string]*ServiceState
}

func (sm *ServiceManager) Init() error {
	if err := sm.initServices(func() []string {
		sl := make([]string, 0)
		for sn, ss := range sm.services {
			if ss.core {
				sl = append(sl, sn)
			}
		}
		return sl
	}()); err != nil {
		return err
	}
	if err := sm.initServices(func() []string {
		sl := make([]string, 0)
		for sn, ss := range sm.services {
			if !ss.core {
				sl = append(sl, sn)
			}
		}
		return sl
	}()); err != nil {
		return err
	}
	for sn, ss := range sm.services {
		ss.service.SetAdditionalDependencies(sm.getServiceAdditionalDependencies(sn))
	}
	return nil
}

func (sm *ServiceManager) Stop() error {
	var wg sync.WaitGroup
	for sn, ss := range sm.services {
		wg.Add(1)
		go func(sn string, ss *ServiceState) {
			c := make(chan interface{})
			go ss.service.UnsetAdditionalDependencies(c)
			<-c
			wg.Done()
		}(sn, ss)
	}
	wg.Wait()
	for {
		running := false
		for sn, ss := range sm.services {
			if !ss.running {
				continue
			}
			running = true
			deprunning := false
			for _, dss := range sm.services {
				if !dss.running {
					continue
				}
				dl, _ := dss.service.Dependencies()
				if funk.ContainsString(dl, sn) {
					deprunning = true
				}
			}
			if !deprunning {
				ss.service.Stop()
				sm.genServiceLogger(sn)("Stopped")
				ss.running = false
			}
		}
		if !running {
			return nil
		}
	}
}

func (sm *ServiceManager) genServiceLogger(sn string) service.LoggerFunc {
	return func(v ...interface{}) {
		log.Println(append([]interface{}{"Service:", sn + ":"}, v...)...)
	}
}

func (sm *ServiceManager) genServiceError(sn string) service.ErrorFunc {
	return func(v ...interface{}) {
		log.Fatal(append([]interface{}{"Service error: ", sn + ": "}, v...)...)
	}
}

func (sm *ServiceManager) initServices(sl []string) error {
	init := 0
	for len(sl) > init {
		for _, sn := range sl {
			ss := sm.services[sn]
			if ss.running {
				continue
			}
			s := ss.service
			dl, _ := ss.service.Dependencies()
			dln := make([]string, 0)
			for _, d := range dl {
				if !sm.services[d].running {
					dln = append(dln, d)
				}
			}
			if len(dln) == 0 {
				err := s.Init(
					sm.getServiceDependencies(sn),
					sm.genServiceLogger(sn),
					sm.genServiceError(sn),
				)
				if err != nil {
					return err
				}
				sm.genServiceLogger(sn)("Started")
				ss.running = true
				init++
			}
		}
	}
	return nil
}

func (sm *ServiceManager) addService(s service.Service) {
	sm.services[s.Name()] = &ServiceState{service: s}
}

func (sm *ServiceManager) addCoreService(s service.Service) {
	sm.services[s.Name()] = &ServiceState{service: s, core: true}
}

func (sm *ServiceManager) getServiceDependencies(sn string) map[string]service.Service {
	dep := make(map[string]service.Service)
	dl, _ := sm.services[sn].service.Dependencies()
	for _, d := range dl {
		dep[d] = sm.services[d].service
	}
	return dep
}

func (sm *ServiceManager) getServiceAdditionalDependencies(sn string) map[string]service.Service {
	dep := make(map[string]service.Service)
	_, dl := sm.services[sn].service.Dependencies()
	for _, d := range dl {
		dep[d] = sm.services[d].service
	}
	return dep
}

func (sm *ServiceManager) checkCircularDependencies() error {
	for sn, ss := range sm.services {
		if err := sm.serviceCheckCircularDependencies(sn, ss.service, &[]string{}); err != nil {
			return err
		}
	}
	return nil
}

func (sm *ServiceManager) serviceCheckCircularDependencies(start string, check service.Service, checked *[]string) error {
	dl, adl := check.Dependencies()
	for _, dn := range append(dl, adl...) {
		if start == dn {
			return errors.New("Circular import detected")
		}
		if !funk.ContainsString(*checked, dn) {
			*checked = append(*checked, dn)
			if err := sm.serviceCheckCircularDependencies(start, sm.services[dn].service, checked); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewServiceManager() (*ServiceManager, error) {
	sm := new(ServiceManager)
	sm.services = make(map[string]*ServiceState)
	for _, s := range services {
		sm.addService(s)
	}
	for _, s := range coreservices {
		sm.addCoreService(s)
	}
	if err := sm.checkCircularDependencies(); err != nil {
		return nil, err
	}
	return sm, nil
}
