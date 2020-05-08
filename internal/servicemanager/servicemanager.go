package servicemanager

import (
	"errors"
	"log"
	"sync"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/ftCommunity/roboheart/internal/services/core/fwver"
	"github.com/ftCommunity/roboheart/internal/services/core/power"
	"github.com/ftCommunity/roboheart/internal/services/core/releasever"
	"github.com/thoas/go-funk"
)

var (
	services = []service.Service{
	}
	coreservices = []service.Service{
		config.Service,
		acm.Service,
		fwver.Service,
		relver.Service,
		power.Service,
	}
)

type ServiceState struct {
	service service.Service
	core    bool
	running bool
	addset  bool //Additional dependencies set
}

type ServiceManager struct {
	services map[string]*ServiceState
}

func (sm *ServiceManager) Init() error {
	//inititialize core services
	if err := sm.initServices(func() []string {
		sl := make([]string, 0)
		//iterate over services
		for sn, ss := range sm.services {
			//filter on core services
			if ss.core {
				sl = append(sl, sn)
			}
		}
		return sl
	}()); err != nil {
		//return in case of error
		return err
	}
	//inititialize non-core services
	if err := sm.initServices(func() []string {
		sl := make([]string, 0)
		//iterate over services
		for sn, ss := range sm.services {
			//filter on non-core services
			if !ss.core {
				sl = append(sl, sn)
			}
		}
		return sl
	}()); err != nil {
		//return in case of error
		return err
	}
	for sn, ss := range sm.services {
		ss.service.SetAdditionalDependencies(sm.getServiceAdditionalDependencies(sn))
	}
	return nil
}

func (sm *ServiceManager) Stop() error {
	var wg sync.WaitGroup
	//iterate over services
	for sn, ss := range sm.services {
		//increment wait group
		wg.Add(1)
		go func(sn string, ss *ServiceState) {
			//unset dependencies
			ss.service.UnsetAdditionalDependencies()
			//decrease wait group
			wg.Done()
		}(sn, ss)
	}
	//wait for all services
	wg.Wait()
	//run forever
	for {
		running := false
		//iterate over all services
		for sn, ss := range sm.services {
			if !ss.running {
				//jump to next service if this is not running
				continue
			}
			//found running service!
			//save this
			running = true
			//assume no dependencies running
			deprunning := false
			//iterate over services to find running depending on this
			for _, dss := range sm.services {
				if !dss.running {
					//jump to next service if this is not running
					continue
				}
				//get dependencies
				dl, _ := dss.service.Dependencies()
				if funk.ContainsString(dl, sn) {
					//save that depending service is running
					deprunning = true
				}
			}
			//initiate stop procedure if no depending service is running
			if !deprunning {
				//stop it
				ss.service.Stop()
				//log that
				sm.genServiceLogger(sn)("Stopped")
				//save that
				ss.running = false
			}
		}
		if !running {
			//end here when all services where already stopped
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

//initServices starts a list of services
//caution: all services must depend on each other or dependcies must be started already
func (sm *ServiceManager) initServices(sl []string) error {
	//save how many are inititialized
	init := 0
	for len(sl) > init {
		//iterate over services
		for _, sn := range sl {
			ss := sm.services[sn]
			if ss.running {
				//ignore if running
				continue
			}
			//get service
			s := ss.service
			//get dependencies
			dl, _ := ss.service.Dependencies()
			dm := false
			//iterate over dependencies
			for _, d := range dl {
				if !sm.services[d].running {
					//set flag if dependency is not running
					dm = true
				}
			}
			if !dm {
				//inititialize if there are no missing dependencies
				err := s.Init(
					sm.getServiceDependencies(sn),
					sm.genServiceLogger(sn),
					sm.genServiceError(sn),
				)
				if err != nil {
					return err
				}
				//log running
				sm.genServiceLogger(sn)("Started")
				//save running state
				ss.running = true
				//increase counter
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
	//iterate over services
	for sn, ss := range sm.services {
		//return error if one service has a circular dependency
		if err := sm.serviceCheckCircularDependencies(sn, ss.service, &[]string{}); err != nil {
			return err
		}
	}
	return nil
}

func (sm *ServiceManager) serviceCheckCircularDependencies(start string, check service.Service, checked *[]string) error {
	//start: start service name
	//check: current service
	//checked: already checked services sincde start

	//get dependencies
	dl, adl := check.Dependencies()
	//range over dependencies and additional dependencies
	for _, dn := range append(dl, adl...) {
		if start == dn {
			//circular dependency if we hit the start again
			return errors.New("Circular import detected")
		}
		if !funk.ContainsString(*checked, dn) {
			//check tree if not in checked
			//save service name
			*checked = append(*checked, dn)
			//run the next check
			if err := sm.serviceCheckCircularDependencies(start, sm.services[dn].service, checked); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sm *ServiceManager) checkDependencyNames() error {
	//iterate over services
	for sn, ss := range sm.services {
		//get dependencies
		dl, adl := ss.service.Dependencies()
		//range over dependencies and additional dependencies
		for _, dn := range append(dl, adl...) {
			if _, ok := sm.services[dn]; !ok {
				//return error if service is not unknown
				return errors.New("Service " + sn + " requests dependency " + dn + " which is not definded")
			}
		}
	}
	return nil
}

func NewServiceManager() (*ServiceManager, error) {
	//create ServiceManager amd inititialize it
	sm := new(ServiceManager)
	sm.services = make(map[string]*ServiceState)
	//add services
	for _, s := range services {
		sm.addService(s)
	}
	for _, s := range coreservices {
		sm.addCoreService(s)
	}
	//run checks
	if err := sm.checkDependencyNames(); err != nil {
		return nil, err
	}
	if err := sm.checkCircularDependencies(); err != nil {
		return nil, err
	}
	return sm, nil
}
