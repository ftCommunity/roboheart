package manifest

import (
	"encoding/json"
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
)

type (
	GetStartup   func(Configurator) []instance.ID
	InitFunc     func(instance.ID, instance.LoggerFunc, instance.ErrorFunc, instance.SelfKillFunc, Configurator) instance.Instance
	ConfigFunc   func(json.RawMessage) (Configurator, error) //argument == nil value => no config
	Configurator interface{}
)

type ServiceManifest struct {
	Name         string
	Instantiable bool
	GetStartup   GetStartup //can be nil
	InitFunc     InitFunc
	ConfigFunc   ConfigFunc
}
