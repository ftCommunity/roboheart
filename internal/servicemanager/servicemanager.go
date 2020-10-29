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
			si.getBase().Stop()
			delete(ss.instances, si.id.Instance)
		}
	}
}

func (sm *ServiceManager) newInstance(id instance.ID) error {
	var ss *ServiceState
	if ss := sm.services[id.Name]; ss == nil {
		return errors.New("Service " + id.Name + " is unknown")
	}
	if //goland:noinspection GoNilness
	err := ss.init(id); err != nil {
		log.Fatal(err)
	}
	//goland:noinspection GoNilness
	si := ss.get(id)
	si.load()
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

func (sm *ServiceManager) genServiceLogger(id instance.ID) instance.LoggerFunc {
	var sn string
	if id.Instance == NON_INSTANCE_NAME {
		sn = id.Name
	} else {
		sn = id.Name + "." + id.Instance
	}
	sn = "\"" + sn + "\""
	return func(v ...interface{}) {
		log.Println(append([]interface{}{"Log from instance", sn + ":"}, v...)...)
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

func (sm *ServiceManager) genServiceError(id instance.ID) instance.ErrorFunc {
	var sn string
	if id.Instance == NON_INSTANCE_NAME {
		sn = id.Name
	} else {
		sn = id.Name + "." + id.Instance
	}
	sn = "\"" + sn + "\""
	return func(v ...interface{}) {
		log.Println(append([]interface{}{"Error on instance", sn + ":"}, v...)...)
		sm.serviceslock.Lock()
		defer sm.serviceslock.Unlock()
		is := sm.get(id)
		if fs := is.instance.forcestop; fs != nil {
			fs.ForceStop()
		}
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
		return errors.New("Service name must not be empty")
	}
	if _, ok := sm.services[m.Name]; ok {
		return errors.New("Service " + m.Name + " loaded twice")
	}
	if m.InitFunc == nil {
		return errors.New("Service " + m.Name + " does not have InitFunc")
	}
	sm.services[m.Name] = newServiceState(m, builtin)
	if gsuf := sm.services[m.Name].GetStartup; gsuf != nil {
		for _, suid := range gsuf(sm.services[m.Name].configurator) {
			if err := sm.newInstance(suid); err != nil {
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

func NewServiceManager() (*ServiceManager, error) {
	//create ServiceManager amd initialize it
	sm := new(ServiceManager)
	sm.services = make(map[string]*ServiceState)
	sm.exposed = newExposed(sm)
	//add services
	for _, m := range services.Services {
		if err := sm.loadService(m, true); err != nil {
			return nil, err
		}
	}
	sm.workercheck = make(chan interface{})
	return sm, nil
}
