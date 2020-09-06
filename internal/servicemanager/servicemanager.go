package servicemanager

import (
	"errors"
	"github.com/ftCommunity-roboheart/roboheart/internal/services"
	"log"
	"strings"
	"sync"

	"github.com/ftCommunity-roboheart/roboheart/package/service"
	"github.com/thoas/go-funk"
)

type ServiceManager struct {
	services map[string]*ServiceState
	exposed  *exposed
}

func (sm *ServiceManager) Init() {
	good := false
	for !good {
		good = true
		for _, ss := range sm.services {
			ok := ss.tryRun()
			if ok {
				if !ss.setReadyAdeps() {
					good = false
				}
			} else {
				good = false
			}
		}
	}
}

func (sm *ServiceManager) Stop() {
	good := false
	for !good {
		good = true
		for _, ss := range sm.services {
			ok := ss.tryStop()
			if !ok {
				good = false
			}
		}
	}
}

func (sm *ServiceManager) emergencyStop() {
	var wg sync.WaitGroup
	for _, ss := range sm.services {
		wg.Add(1)
		go func(s *ServiceState) {
			s.emerstop()
			wg.Done()
		}(ss)
	}
	wg.Wait()
}

func (sm *ServiceManager) genServiceLogger(sn string) service.LoggerFunc {
	return func(v ...interface{}) {
		log.Println(append([]interface{}{"Service:", sn + ":"}, v...)...)
	}
}

func (sm *ServiceManager) genServiceError(sn string) service.ErrorFunc {
	return func(v ...interface{}) {
		log.Println(append([]interface{}{"Service error:", sn + ":"}, v...)...)
		log.Println("Emergency stop for all services")
		sm.emergencyStop()
		log.Fatal("Stopped after previous error in service " + sn)
	}
}

func (sm *ServiceManager) getServiceDependencies(sn string) map[string]service.Service {
	dep := make(map[string]service.Service)
	if ds := sm.services[sn].service.depending; ds != nil {
		for _, d := range sm.services[sn].deps.Deps {
			dep[d] = sm.services[d].service.base
		}
	}
	return dep
}

func (sm *ServiceManager) checkCircularDependencies() error {
	//iterate over services
	for sn, ss := range sm.services {
		//return error if one service has a circular dependency
		if ds := ss.service.depending; ds != nil {
			if err := sm.serviceCheckCircularDependencies(sn, ds, &[]string{}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sm *ServiceManager) serviceCheckCircularDependencies(start string, check service.DependingService, checked *[]string) error {
	//start: start service name
	//check: current service
	//checked: already checked services since start

	//get dependencies
	deps := check.Dependencies()
	dl, adl := deps.Deps, deps.ADeps
	//range over dependencies and additional dependencies
	for _, dn := range append(dl, adl...) {
		if start == dn {
			//circular dependency if we hit the start again
			return errors.New("circular import detected")
		}
		if !funk.ContainsString(*checked, dn) {
			//check tree if not in checked
			//save service name
			*checked = append(*checked, dn)
			//run the next check
			if nds := sm.services[dn].service.depending; nds != nil {
				if err := sm.serviceCheckCircularDependencies(start, nds, checked); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (sm *ServiceManager) checkDependencyNames() error {
	var errs []string
	//iterate over services
	for sn, ss := range sm.services {
		var deps []string
		if ds := ss.service.depending; ds != nil {
			//get dependencies
			//range over dependencies and additional dependencies
			for _, dn := range append(ss.deps.Deps, ss.deps.ADeps...) {
				if _, ok := sm.services[dn]; !ok {
					deps = append(deps, dn)
				}
			}
		}
		if len(deps) > 0 {
			errs = append(errs, sn+"->("+strings.Join(deps, ",")+")")
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(append([]string{"Service dependency error(s):"}, errs...), "\n"))
	}
	return nil
}

func NewServiceManager() (*ServiceManager, error) {
	//create ServiceManager amd initialize it
	sm := new(ServiceManager)
	sm.services = make(map[string]*ServiceState)
	//add services
	for _, s := range services.Services {
		ss := newServiceStateBuiltin(sm, s)
		sm.services[ss.name] = ss
	}
	//run checks
	if err := sm.checkDependencyNames(); err != nil {
		return nil, err
	}
	if err := sm.checkCircularDependencies(); err != nil {
		return nil, err
	}
	//load second stage
	for _, ss := range sm.services {
		ss.loadDepData()
	}
	sm.exposed = newExposed(sm)
	return sm, nil
}
