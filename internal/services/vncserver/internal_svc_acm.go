package vncserver

import (
	"errors"

	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (v *vncserver) initSvcAcm(services servicehelpers.ServiceList) error {
	var ok bool
	v.acm, ok = services["acm"].(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	v.acm.RegisterPermission(PERMISSION, map[string]bool{"user": true, "app": false}, map[string]string{})
	return nil
}
