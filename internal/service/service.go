package service

type Service interface {
	Init(map[string]Service, LoggerFunc, ErrorFunc)
	Name() string
	Stop()
}

type EmergencyStoppableService interface {
	Service
	EmergencyStop()
}

type DependingService interface {
	Service
	Dependencies() ServiceDependencies
}

type AddDependingService interface {
	DependingService
	SetAdditionalDependencies(map[string]Service)
	UnsetAdditionalDependencies([]string)
}

type ManagingService interface {
	SetServiceManager(ServiceManager)
}

type ServiceManager interface {
}

type LoggerFunc func(...interface{})
type ErrorFunc func(...interface{})

type ServiceDependencies struct {
	Deps, ADeps []string
}
