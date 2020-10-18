package manifest

import "github.com/ftCommunity-roboheart/roboheart/package/instance"

type (
	GetStartup func() []instance.ID
	InitFunc   func(instance.ID, instance.LoggerFunc, instance.ErrorFunc, instance.SelfKillFunc) instance.Instance
)

type ServiceManifest struct {
	Name         string
	Instantiable bool
	GetStartup   GetStartup
	InitFunc     InitFunc
}
