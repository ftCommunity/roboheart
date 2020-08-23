package config

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/robotime"
)

func (c *config) initSvcRoboTime(svc service.Service) error {
	var ok bool
	c.rt, ok = svc.(robotime.RoboTime)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
