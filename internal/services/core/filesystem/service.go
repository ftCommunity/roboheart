package filesystem

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/spf13/afero"
)

type filesystem struct {
	afero.Fs
}

func (f *filesystem) Init(map[string]service.Service, service.LoggerFunc, service.ErrorFunc) {
	f.Fs = afero.NewOsFs()
}

func (f *filesystem) Stop() {}

func (f *filesystem) Name() string { return "filesystem" }
