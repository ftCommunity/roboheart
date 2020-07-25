package service

type Service interface {
	Init(map[string]Service, LoggerFunc, ErrorFunc) error
	Name() string
}

type StoppableService interface {
	Service
	Stop() error
}

type EmergencyStoppableService interface {
	StoppableService
	EmergencyStop()
}

type DependingService interface {
	Service
	Dependencies() ([]string, []string)
}

type AddDependingService interface {
	DependingService
	SetAdditionalDependencies(map[string]Service) error
}

type AddDependingUnsetService interface {
	AddDependingService
	UnsetAdditionalDependencies()
}

type LoggerFunc func(...interface{})
type ErrorFunc func(...interface{})
