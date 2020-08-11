package servicehelpers

import (
	"errors"
	"strings"

	"github.com/ftCommunity/roboheart/internal/service"
)

func checkDependencies(req []string, svcs map[string]service.Service) error {
	missing := make([]string, 0)
	for _, s := range req {
		if _, f := svcs[s]; !f {
			missing = append(missing, s)
		}
	}
	if len(missing) > 0 {
		return errors.New("Missing service(s): " + strings.Join(missing, ", "))
	}
	return nil
}

func CheckMainDependencies(ds service.DependingService, svcs map[string]service.Service) error {
	sd := ds.Dependencies()
	return checkDependencies(sd.Deps, svcs)
}

type ServiceList map[string]service.Service
type ServiceInitializers map[string]func(service.Service) error
type AdditionalServiceInitializers map[string]func(service.Service)

func InitializeDependencies(sl ServiceList, dcl ServiceInitializers) error {
	for sn, s := range sl {
		if si, ok := dcl[sn]; ok {
			if err := si(s); err != nil {
				return err
			}
		}
	}
	return nil
}

func InitializeAdditionalDependencies(sl ServiceList, dcl AdditionalServiceInitializers) {
	for sn, s := range sl {
		if si, ok := dcl[sn]; ok {
			si(s)
		}
	}
}
