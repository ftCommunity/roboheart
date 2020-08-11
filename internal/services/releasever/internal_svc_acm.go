package relver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
)

func (r *relver) initSvcAcm(svc service.Service) error {
	var ok bool
	r.acm, ok = svc.(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	r.acm.RegisterPermission(PERMISSION_UPDATE, map[string]bool{"user": true, "app": false}, map[string]string{})
	return nil
}
