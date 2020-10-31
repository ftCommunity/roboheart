package manifest

import (
	"encoding/json"
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
)

type (
	Configurator interface{}
)

type (
	ConfigLoaderFunc        func(json.RawMessage) (Configurator, error) //argument == nil value => no config
	ServiceLoaderFunc       func(Configurator) error
	GetStartupInstancesFunc func(Configurator) []instance.ID
	InstanceInitFunc        func(instance.ID, instance.LoggerFunc, instance.ErrorFunc, instance.SelfKillFunc, Configurator) instance.Instance
)

type ServiceManifest struct {
	Name                    string
	Instantiable            bool
	ConfigLoaderFunc        ConfigLoaderFunc        //can be nil
	ServiceLoaderFunc       ServiceLoaderFunc       //can be nil
	GetStartupInstancesFunc GetStartupInstancesFunc //can be nil
	InstanceInitFunc        InstanceInitFunc
}
