package servicemanager

import (
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
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
		managing  instance.ManagingInstance
	}
	running  bool
	startup  bool // saves whether this was requested on startup
	lastrdep time.Time
	deps     struct {
		deps, rdeps *instance.Dependencies
	}
	logger     instance.LoggerFunc
	error      instance.ErrorFunc
	selfkiller instance.SelfKillFunc
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
	if mi, ok := is.getBase().(instance.ManagingInstance); ok && is.ss.builtin {
		is.instance.managing = mi
	}
}

func (is *InstanceState) load() {
	is.deps.deps = new(instance.Dependencies)
	is.deps.rdeps = new(instance.Dependencies)
	is.instance.base = is.ss.ServiceManifest.InitFunc(is.id)
	is.loadInterfaces()
	is.logger = is.sm.genServiceLogger(is.id)
	is.error = is.sm.genServiceError(is.id)
	is.selfkiller = is.sm.genSelfKillFunc(is.id)
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
			if is.sm.newInstance(nd) != nil {
				continue newdeps
			}
			ni = is.sm.services[nd.Name].get(nd)
		}
		if !ni.running {
			continue newdeps
		}
		is.deps.deps.Add(nd)
		ni.setRdep(is.id)
		di.SetDependency(ni.instance.base)
	}
}
