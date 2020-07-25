package deviceinfo

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (d *deviceinfo) initSvcFileSystem(services servicehelpers.ServiceList) error {
	var ok bool
	d.fs, ok = services["filesystem"].(filesystem.FileSystem)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
