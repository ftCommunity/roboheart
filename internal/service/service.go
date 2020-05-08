package service

type Service interface {
	Init(map[string]Service, LoggerFunc, ErrorFunc) error
	Stop() error
	Name() string
	Dependencies() ([]string, []string)
	SetAdditionalDependencies(map[string]Service) error
	UnsetAdditionalDependencies()
}

type LoggerFunc func(...interface{})
type ErrorFunc func(...interface{})
