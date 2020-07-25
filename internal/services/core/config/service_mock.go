package config

import (
	"github.com/digineo/go-uci"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/threadmanager"
	"github.com/spf13/afero"
)

func NewMock(logger service.LoggerFunc, e service.ErrorFunc) (Config, error) {
	c := &config{}
	c.logger = logger
	c.error = e
	c.tree = uci.NewTreeFromFs(afero.NewMemMapFs())
	c.tm = threadmanager.NewThreadManager(c.logger, c.error)
	c.tm.Load("commit", c.configCommitThread)
	c.tm.Start("commit")
	return c, nil
}
