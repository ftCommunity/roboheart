package servicemanager

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/thoas/go-funk"
)

type ServiceState struct {
	service struct {
		base         service.Service
		stoppable    service.StoppableService
		depending    service.DependingService
		adddepending service.AddDependingService
		addunset     service.AddDependingUnsetService
	}
	running bool
	addset  bool //Additional dependencies set
	logger  service.LoggerFunc
	error   service.ErrorFunc
}

type ServiceManager struct {
	services map[string]*ServiceState
}

func (sm *ServiceManager) Init() error {
	//inititialize services
	if err := sm.initServices(func() []string {
		sl := make([]string, 0)
		//iterate over services
		for sn := range sm.services {
			sl = append(sl, sn)
		}
		return sl
	}()); err != nil {
		//return in case of error
		return err
	}
	for sn, ss := range sm.services {
		if ads := ss.service.adddepending; ads != nil {
			ads.SetAdditionalDependencies(sm.getServiceAdditionalDependencies(sn))
		}
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
			if aus := ss.service.addunset; aus != nil {
				aus.UnsetAdditionalDependencies()
			}
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
				if ds := dss.service.depending; ds != nil {
					dl, _ := ds.Dependencies()
					if funk.ContainsString(dl, sn) {
						//save that depending service is running
						deprunning = true
					}
				}
			}
			//initiate stop procedure if no depending service is running
			if !deprunning {
				//stop it
				if sts := ss.service.stoppable; sts != nil {
					sts.Stop()
					//log that
					ss.logger("Stopped")
				} else {
					ss.logger("Marked as stopped")
				}
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
			dm := false
			//resolve dependencies
			if ds := ss.service.depending; ds != nil {
				dl, _ := ds.Dependencies()
				//iterate over dependencies
				for _, d := range dl {
					if !sm.services[d].running {
						//set flag if dependency is not running
						dm = true
					}
				}
			}
			if !dm {
				//inititialize if there are no missing dependencies
				err := ss.service.base.Init(
					sm.getServiceDependencies(sn),
					ss.logger,
					ss.error,
				)
				if err != nil {
					return err
				}
				//log running
				ss.logger("Started")
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
	sn := s.Name()
	ss := &ServiceState{}
	ss.service.base = s
	if sts, ok := s.(service.StoppableService); ok {
		ss.service.stoppable = sts
	}
	if ds, ok := s.(service.DependingService); ok {
		ss.service.depending = ds
	}
	if ads, ok := s.(service.AddDependingService); ok {
		ss.service.adddepending = ads
	}
	if aus, ok := s.(service.AddDependingUnsetService); ok {
		ss.service.addunset = aus
	}
	ss.logger = sm.genServiceLogger(sn)
	ss.error = sm.genServiceError(sn)
	sm.services[sn] = ss
}

func (sm *ServiceManager) getServiceDependencies(sn string) map[string]service.Service {
	dep := make(map[string]service.Service)
	if ds := sm.services[sn].service.depending; ds != nil {
		dl, _ := ds.Dependencies()
		for _, d := range dl {
			dep[d] = sm.services[d].service.base
		}
	}
	return dep
}

func (sm *ServiceManager) getServiceAdditionalDependencies(sn string) map[string]service.Service {
	dep := make(map[string]service.Service)
	if ds := sm.services[sn].service.depending; ds != nil {
		_, dl := ds.Dependencies()
		for _, d := range dl {
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
	errs := []string{}
	//iterate over services
	for sn, ss := range sm.services {
		deps := []string{}
		if ds := ss.service.depending; ds != nil {
			//get dependencies
			dl, adl := ds.Dependencies()
			//range over dependencies and additional dependencies
			for _, dn := range append(dl, adl...) {
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
	//create ServiceManager amd inititialize it
	sm := new(ServiceManager)
	sm.services = make(map[string]*ServiceState)
	//add services
	for _, s := range getServices() {
		sm.addService(s)
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
