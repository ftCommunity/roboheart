package relver

import (
	"errors"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (r *relver) initSvcAcm(services servicehelpers.ServiceList) error {
	var ok bool
	r.acm, ok = services["acm"].(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	r.acm.RegisterPermission(PERMISSION_UPDATE, map[string]bool{"user": true, "app": false}, map[string]string{})
	return nil
}
