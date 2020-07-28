package servicehelpers

import "github.com/ftCommunity/roboheart/internal/service"

type dummyService struct {
	deps, adeps []string
}

func (dummyService) Init(map[string]service.Service, service.LoggerFunc, service.ErrorFunc) error {
	return nil
}
func (dummyService) Name() string {
	return ""
}

func (d dummyService) Dependencies() ([]string, []string) {
	return d.deps, d.adeps
}
