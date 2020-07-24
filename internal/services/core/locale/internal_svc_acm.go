package locale

import (
	"errors"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (l *locale) initSvcAcm(services servicehelpers.ServiceList) error {
	var ok bool
	l.acm, ok = services["acm"].(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	l.acm.RegisterPermission(PERMISSION, map[string]bool{"user": true, "app": false}, map[string]string{})
	return nil
}
