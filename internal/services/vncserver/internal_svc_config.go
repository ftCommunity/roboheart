package vncserver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"

	"github.com/ftCommunity/roboheart/internal/services/core/config"
)

func (v *vncserver) initSvcConfig(svc service.Service) error {
	var ok bool
	v.config, ok = svc.(config.Config)
	if !ok {
		return errors.New("Type assertion error")
	}
	v.sconfig = v.config.GetServiceConfig(v)
	if err := v.sconfig.AddSection(CONFIG_SECTION, CONFIG_TYPE); err != nil {
		return err
	}
	return nil
}
