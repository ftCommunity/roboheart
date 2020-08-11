package config

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
)

func (c *config) initSvcFileSystem(svc service.Service) error {
	var ok bool
	c.fs, ok = svc.(filesystem.FileSystem)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
