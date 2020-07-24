package vncserver

import (
	"errors"

	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (v *vncserver) initSvcConfig(services servicehelpers.ServiceList) error {
	var ok bool
	v.config, ok = services["config"].(config.Config)
	if !ok {
		return errors.New("Type assertion error")
	}
	v.sconfig = v.config.GetServiceConfig(v)
	if err := v.sconfig.AddSection(CONFIG_SECTION, CONFIG_TYPE); err != nil {
		return err
	}
	return nil
}
