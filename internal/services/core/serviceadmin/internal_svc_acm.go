package serviceadmin

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
)

func (s *serviceadmin) initSvcAcm(svc service.Service) error {
	var ok bool
	s.acm, ok = svc.(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
