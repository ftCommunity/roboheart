package filesystem

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/spf13/afero"
)

func NewMock(service.LoggerFunc, service.ErrorFunc) (FileSystem, error) {
	return afero.NewMemMapFs(), nil
}
