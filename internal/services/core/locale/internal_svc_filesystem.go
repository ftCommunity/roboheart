package locale

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (l *locale) initSvcFileSystem(services servicehelpers.ServiceList) error {
	var ok bool
	l.fs, ok = services["filesystem"].(filesystem.FileSystem)
	if !ok {
		return errors.New("Type assertion error")
	}
	return nil
}
