package servicemanager

import (
	"errors"
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
	"github.com/ftCommunity-roboheart/roboheart/package/manifest"
)

type ServiceState struct {
	manifest.ServiceManifest
	builtin   bool
	instances map[string]*InstanceState
	sm        *ServiceManager
}

func (ss *ServiceState) init(id instance.ID) error {
	if _, ok := ss.instances[id.Instance]; ok {
		return errors.New("Instance " + id.Instance + " does already exist")
	}
	if id.Instance != NON_INSTANCE_NAME && !ss.ServiceManifest.Instantiation {
		return errors.New("Service " + ss.ServiceManifest.Name + " cannot be instantiated")
	}
	is := &InstanceState{}
	is.sm = ss.sm
	is.ss = ss
	is.id = id
	ss.instances[id.Instance] = is
	return nil
}

func (ss *ServiceState) get(id instance.ID) *InstanceState {
	if si, ok := ss.instances[id.Instance]; ok {
		return si
	} else {
		return nil
	}
}

func newServiceState(m manifest.ServiceManifest, builtin bool) *ServiceState {
	ss := new(ServiceState)
	ss.ServiceManifest = m
	ss.instances = make(map[string]*InstanceState)
	ss.builtin = builtin
	return ss
}
