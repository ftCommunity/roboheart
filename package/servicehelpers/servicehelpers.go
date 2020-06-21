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
	req, _ := ds.Dependencies()
	return checkDependencies(req, svcs)
}

func CheckAdditionalDependencies(ds service.DependingService, svcs map[string]service.Service) error {
	_, req := ds.Dependencies()
	return checkDependencies(req, svcs)
}
