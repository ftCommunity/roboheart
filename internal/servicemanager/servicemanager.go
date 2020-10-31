package servicemanager

import (
	"errors"
	"github.com/ftCommunity-roboheart/roboheart/internal/services"
	"github.com/ftCommunity-roboheart/roboheart/package/manifest"
	"github.com/thoas/go-funk"
	"log"
	"plugin"
	"sync"

	"github.com/ftCommunity-roboheart/roboheart/package/instance"
)

type ServiceManager struct {
	services       map[string]*ServiceState
	exposed        *exposed
	wg             sync.WaitGroup
	workercall     bool
	workerstop     bool
	workercheck    chan interface{}
	workercalllock sync.Mutex
	serviceslock   sync.Mutex
	config         *config
}

func (sm *ServiceManager) Init() {
	sm.wg.Add(1)
	go sm.worker()
	go sm.triggerWorker()
}

func (sm *ServiceManager) Stop() {
	sm.stopWorker()
	sm.wg.Wait()
	// we are missing many steps here...
	// but!
	// we do not need them!
	// after the action everything is clear!
	for _, ss := range sm.services {
		for _, si := range ss.instances {
			for _, rd := range *si.deps.rdeps {
				sm.get(rd).instance.depending.UnsetDependency(si.id)
			}
			si.stop()
			delete(ss.instances, si.id.Instance)
		}
	}
}

func (sm *ServiceManager) newInstance(id instance.ID, startup bool) error {
	ss := sm.services[id.Name]
	if ss == nil {
		return errors.New("Service " + id.Name + " is unknown")
	}
	if err := ss.init(id); err != nil {
		log.Fatal(err)
	}
	si := ss.get(id)
	if err := si.load(); err != nil {
		return err
	}
	si.startup = startup
	if mi := si.instance.managing; mi != nil {
		mi.SetServiceManager(sm.exposed)
	}
	if di := si.instance.depending; di != nil {
		di.SetServiceListGetter(sm.getServiceList)
		di.SetDependenciesChangedHandler(si.updateDependencies)
	}
	return nil
}

func (sm *ServiceManager) getServiceList() []string {
	return funk.Keys(sm.services).([]string)
}

func (sm *ServiceManager) genServiceLogger(idstr string) instance.LoggerFunc {
	return func(v ...interface{}) {
		log.Println(append([]interface{}{"Log from instance", "\"" + idstr + "\"" + ":"}, v...)...)
	}
}

func (sm *ServiceManager) genSelfKillFunc(id instance.ID) instance.SelfKillFunc {
	return func() {
		sm.serviceslock.Lock()
		defer sm.serviceslock.Unlock()
		si := sm.get(id)
		for _, di := range *si.deps.deps {
			si.instance.depending.UnsetDependency(di)
			sm.get(di).unsetRdep(id)
		}
		for _, di := range *si.deps.rdeps {
			sm.get(di).instance.depending.UnsetDependency(id)
		}
		delete(sm.services[id.Name].instances, id.Instance)
	}
}

func (sm *ServiceManager) genServiceError(id instance.ID, idstr string) instance.ErrorFunc {
	return func(v ...interface{}) {
		log.Println(append([]interface{}{"Error on instance", "\"" + idstr + "\"" + ":"}, v...)...)
		sm.serviceslock.Lock()
		defer sm.serviceslock.Unlock()
		is := sm.get(id)
		is.forcestop()
		is.running = false
		for _, rd := range *is.deps.rdeps {
			sm.get(rd).instance.depending.UnsetDependency(id)
			is.unsetRdep(rd)
			sm.get(rd).deps.deps.Delete(id)
		}
		for _, d := range *is.deps.deps {
			// no unset needed because instance is stopped
			is.deps.deps.Delete(d)
			sm.get(d).unsetRdep(id)
		}
	}
}

func (sm *ServiceManager) get(id instance.ID) *InstanceState {
	if ss, ok := sm.services[id.Name]; ok {
		return ss.get(id)
	} else {
		return nil
	}
}

func (sm *ServiceManager) forAllInstances(f func(is *InstanceState) error) error {
	for _, ss := range sm.services {
		for _, is := range ss.instances {
			if err := f(is); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sm *ServiceManager) loadFromPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}
	s, err := p.Lookup(PLUGIN_MANIFEST_SYMBOL)
	if err != nil {
		return err
	}
	m, ok := s.(manifest.ServiceManifest)
	if !ok {
		return errors.New("error reading manifest")
	}
	if err = sm.loadService(m, false); err != nil {
		return err
	}
	return nil
}

func (sm *ServiceManager) loadService(m manifest.ServiceManifest, builtin bool) error {
	sm.serviceslock.Lock()
	defer sm.serviceslock.Unlock()
	if m.Name == "" {
		return errors.New("service name must not be empty")
	}
	if _, ok := sm.services[m.Name]; ok {
		return errors.New("service " + m.Name + " loaded twice")
	}
	if m.InstanceInitFunc == nil {
		return errors.New("service " + m.Name + " does not have InitFunc")
	}
	ss, err := newServiceState(m, builtin, sm)
	if err != nil {
		return err
	}
	sm.services[m.Name] = ss
	if gsuf := ss.GetStartupInstancesFunc; gsuf != nil {
		for _, suid := range gsuf(ss.configurator) {
			if err := sm.newInstance(suid, true); err != nil {
				return err
			}
		}
	}
	return sm.forAllInstances(func(is *InstanceState) error {
		if di := is.instance.depending; di != nil {
			di.OnServiceListChanged()
		}
		return nil
	})
}

func NewServiceManager(config []byte) (*ServiceManager, error) {
	//create ServiceManager amd initialize it
	sm := new(ServiceManager)
	sm.services = make(map[string]*ServiceState)
	sm.exposed = newExposed(sm)
	if c, err := readConfig(config); err != nil {
		return nil, err
	} else {
		sm.config = c
	}
	//add services
	for _, m := range func() []manifest.ServiceManifest {
		var ml []manifest.ServiceManifest
		for _, sp := range services.ServiceProviders {
			ml = append(ml, sp...)
		}
		return ml
	}() {
		if err := sm.loadService(m, true); err != nil {
			return nil, err
		}
	}
	sm.workercheck = make(chan interface{})
	return sm, nil
}
