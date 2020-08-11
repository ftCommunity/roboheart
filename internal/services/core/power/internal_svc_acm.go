package power

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
)

func (p *power) initSvcAcm(svc service.Service) error {
	var ok bool
	p.acm, ok = svc.(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	p.acm.RegisterPermission(PERMISSION, map[string]bool{"user": true, "app": false}, map[string]string{})
	return nil
}
