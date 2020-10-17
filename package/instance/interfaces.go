package instance

type Instance interface {
	Init(LoggerFunc, ErrorFunc)
	Stop()
	ID() ID
}

type ForceStoppableInstance interface {
	ForceStop() //this should just stop threads and clean up fast
}

type DependingInstance interface {
	Instance
	Dependencies() Dependencies
	SetDependency(Instance)
	UnsetDependency(ID)
	OnServiceListChanged([]string)
	SetServiceListGetter(func() []string)
	SetDependenciesChangedHandler(func())
}

type ManagingInstance interface {
	Instance
	SetServiceManager(ServiceManager)
}

type ServiceManager interface {
	//todo
}
