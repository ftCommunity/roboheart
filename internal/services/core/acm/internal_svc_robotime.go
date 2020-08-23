package acm

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/robotime"
)

func (a*acm) initSvcRoboTime(svc service.Service) error {
	var ok bool
	a.rt, ok = svc.(robotime.RoboTime)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
