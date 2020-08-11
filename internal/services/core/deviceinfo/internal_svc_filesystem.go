package deviceinfo

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
)

func (d *deviceinfo) initSvcFileSystem(svc service.Service) error {
	var ok bool
	d.fs, ok = svc.(filesystem.FileSystem)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
