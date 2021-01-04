package servicemanager

import (
	"errors"
	"github.com/servicemngr/core/package/instance"
	"github.com/servicemngr/core/package/manifest"
)

type ServiceState struct {
	manifest.ServiceManifest
	builtin      bool
	instances    map[string]*InstanceState
	configurator manifest.Configurator
	sm           *ServiceManager
}

func (ss *ServiceState) init(id instance.ID) error {
	if id.Name != ss.ServiceManifest.Name {
		return errors.New("Service " + ss.ServiceManifest.Name + ": name mismatch: " + id.Name)
	}
	if _, ok := ss.instances[id.Instance]; ok {
		return errors.New("Instance " + id.Instance + " does already exist")
	}
	if id.Instance != instance.NON_INSTANCE_NAME && !ss.ServiceManifest.Instantiable {
		return errors.New("Service " + ss.ServiceManifest.Name + " cannot be instantiated")
	}
	if id.Instance == instance.NON_INSTANCE_NAME && ss.ServiceManifest.Instantiable {
		return errors.New("Service " + ss.ServiceManifest.Name + " cannot have instance name \"" + instance.NON_INSTANCE_NAME + "\"")
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

func (ss *ServiceState) loadConfig() error {
	if cf := ss.ConfigLoaderFunc; cf != nil {
		if c, err := cf(ss.sm.config.Services[ss.Name]); err != nil {
			return err
		} else {
			ss.configurator = c
		}
	}
	return nil
}

func (ss *ServiceState) loadService() error {
	if sf := ss.ServiceLoaderFunc; sf != nil {
		if err := sf(ss.configurator); err != nil {
			return err
		}
	}
	return nil
}

func newServiceState(m manifest.ServiceManifest, builtin bool, sm *ServiceManager) (*ServiceState, error) {
	ss := new(ServiceState)
	ss.ServiceManifest = m
	ss.instances = make(map[string]*InstanceState)
	ss.builtin = builtin
	ss.sm = sm
	if err := ss.loadConfig(); err != nil {
		return nil, err
	}
	if err := ss.loadService(); err != nil {
		return nil, err
	}
	return ss, nil
}
