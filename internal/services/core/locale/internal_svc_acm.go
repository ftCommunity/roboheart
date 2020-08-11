package locale

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
)

func (l *locale) initSvcAcm(svc service.Service) error {
	var ok bool
	l.acm, ok = svc.(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	l.acm.RegisterPermission(PERMISSION, map[string]bool{"user": true, "app": false}, map[string]string{})
	return nil
}
