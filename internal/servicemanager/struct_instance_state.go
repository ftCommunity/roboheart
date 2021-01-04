package servicemanager

import (
	"errors"
	"github.com/servicemngr/core/package/instance"
	"log"
	"time"
)

type InstanceState struct {
	sm       *ServiceManager
	ss       *ServiceState
	id       instance.ID
	instance struct {
		base      instance.Instance
		forcestop instance.ForceStoppableInstance
		depending instance.DependingInstance
		dependent instance.DependentInstance
		managing  instance.ManagingInstance
	}
	running  bool
	startup  bool // saves whether this was requested on startup
	created  time.Time
	lastrdep time.Time
	deps     struct {
		deps, rdeps *instance.Dependencies
	}
}

func (is *InstanceState) getBase() instance.Instance {
	return is.instance.base
}

func (is *InstanceState) loadInterfaces() {
	if fsi, ok := is.getBase().(instance.ForceStoppableInstance); ok {
		is.instance.forcestop = fsi
	}
	if di, ok := is.getBase().(instance.DependingInstance); ok {
		is.instance.depending = di
	}
	if di, ok := is.getBase().(instance.DependentInstance); ok {
		is.instance.dependent = di
	}
	if mi, ok := is.getBase().(instance.ManagingInstance); ok && is.ss.builtin {
		is.instance.managing = mi
	}
}

func (is *InstanceState) load() error {
	is.deps.deps = new(instance.Dependencies)
	is.deps.rdeps = new(instance.Dependencies)
	is.instance.base = is.ss.ServiceManifest.InstanceInitFunc(is.id, is.sm.genServiceLogger(is.id.String()), is.sm.genServiceError(is.id, is.id.String()), is.sm.genSelfKillFunc(is.id), is.ss.configurator)
	if is.instance.base == nil {
		return errors.New("InstanceInitFunc returned nil")
	}
	is.loadInterfaces()
	is.created = time.Now()
	return nil
}

func (is *InstanceState) setRdep(id instance.ID) {
	is.deps.rdeps.Add(id)
}

func (is *InstanceState) unsetRdep(id instance.ID) {
	is.deps.rdeps.Delete(id)
	is.lastrdep = time.Now()
}

func (is *InstanceState) updateDependencies() {
	if is.instance.depending == nil || !is.running {
		return
	}
	di := is.instance.depending
	ndeps := di.Dependencies()
	o, n, _ := is.deps.deps.Compare(ndeps)
	for _, od := range o {
		di.UnsetDependency(od)
		is.deps.deps.Delete(od)
		is.sm.services[od.Name].get(od).unsetRdep(is.id)
	}
newdeps:
	for _, nd := range n {
		if _, ok := is.sm.services[nd.Name]; !ok {
			continue newdeps
		}
		var ni *InstanceState
		if ni = is.sm.services[nd.Name].get(nd); ni == nil {
			if is.sm.newInstance(nd, false) != nil {
				continue newdeps
			}
			go is.sm.triggerWorker()
			continue newdeps
		}
		if !ni.running {
			continue newdeps
		}
		if ni.instance.dependent == nil {
			continue newdeps
		}
		is.deps.deps.Add(nd)
		ni.setRdep(is.id)
		go di.SetDependency(ni.instance.dependent.GetDependentInstance(is.id))
	}
}

func (is *InstanceState) start() {
	log.Println("Starting instance \"" + is.id.String() + "\"")
	is.getBase().Start()
}

func (is *InstanceState) stop() {
	log.Println("Stopping instance \"" + is.id.String() + "\"")
	is.getBase().Stop()
}

func (is *InstanceState) forcestop() {
	log.Println("Force-stopping instance \"" + is.id.String() + "\"")
	if fs := is.instance.forcestop; fs != nil {
		fs.ForceStop()
	}
}
