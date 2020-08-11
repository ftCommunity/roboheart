package vncserver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
)

func (v *vncserver) initSvcAcm(svc service.Service) error {
	var ok bool
	v.acm, ok = svc.(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	v.acm.RegisterPermission(PERMISSION, map[string]bool{"user": true, "app": false}, map[string]string{})
	return nil
}
