package instance

type Instance interface {
	Start()
	Stop()
	ID() ID
}

type DependencyInstance interface {
	ID() ID
}

type ForceStoppableInstance interface {
	ForceStop() //this should just stop threads and clean up fast
}

type DependingInstance interface {
	Dependencies() Dependencies
	SetDependency(DependencyInstance)
	UnsetDependency(ID)
	OnServiceListChanged()
	// functions below will only be called on instance creation and will not be called again
	SetServiceListGetter(func() []string)
	SetDependenciesChangedHandler(func())
}

type DependentInstance interface {
	GetDependentInstance(ID) DependencyInstance
}

type ManagingInstance interface {
	SetServiceManager(ServiceManager)
}

type ServiceManager interface {
	LoadFromPlugin(path string) error
}
