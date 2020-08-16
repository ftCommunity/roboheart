package servicemanager

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/servicemanager/exposedstructs"
	"github.com/thoas/go-funk"
)

//this must implement service.ServiceManager
type exposed struct {
	sm *ServiceManager
}

func newExposed(sm *ServiceManager) *exposed {
	e := &exposed{}
	e.sm = sm
	return e
}

func (e *exposed) GetServiceList() []string {
	var sl []string
	for sn := range e.sm.services {
		sl = append(sl, sn)
	}
	return sl
}

func (e *exposed) GetServiceInfo(sn string) (exposedstructs.ServiceInfo, error) {
	ss, ok := e.sm.services[sn]
	if !ok {
		return exposedstructs.ServiceInfo{}, errors.New("Service " + sn + " not found")
	}
	return exposedstructs.ServiceInfo{
		Name:    ss.name,
		Running: ss.running,
		Dependencies: exposedstructs.Dependencies{
			Dependencies:    ss.deps.Deps,
			AddDependencies: ss.deps.ADeps,
			AddDependenciesSet: func() []string {
				var ads []string
				for dsn, ds := range ss.deps.adeps {
					if ds.set {
						ads = append(ads, dsn)
					}
				}
				return ads
			}(),
			ReverseDependencies:           funk.Keys(ss.deps.rdeps).([]string),
			ReverseAdditionalDependencies: funk.Keys(ss.deps.radeps).([]string),
			ReverseAdditionalDependenciesSet: func() []string {
				var ads []string
				for dsn, ds := range ss.deps.radeps {
					if ds.set {
						ads = append(ads, dsn)
					}
				}
				return ads
			}(),
		},
		Builtin: ss.builtin,
	}, nil
}

func (e *exposed) GetServicesInfo() map[string]exposedstructs.ServiceInfo {
	si := make(map[string]exposedstructs.ServiceInfo)
	for sn := range e.sm.services {
		si[sn], _ = e.GetServiceInfo(sn)
	}
	return si
}
