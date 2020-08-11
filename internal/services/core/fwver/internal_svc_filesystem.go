package fwver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
)

func (f *fwver) initSvcFileSystem(svc service.Service) error {
	var ok bool
	f.fs, ok = svc.(filesystem.FileSystem)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
