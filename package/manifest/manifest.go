package manifest

import "github.com/ftCommunity-roboheart/roboheart/package/instance"

type (
	GetStartup func() []instance.ID
	InitFunc   func(id instance.ID) instance.Instance
)

type ServiceManifest struct {
	Name          string
	Instantiation bool
	GetStartup    GetStartup
	InitFunc      InitFunc
}
