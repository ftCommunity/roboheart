package config

import (
	"github.com/digineo/go-uci"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/ftCommunity/roboheart/package/threadmanager"
	"github.com/spf13/afero"
)

type config struct {
	logger service.LoggerFunc
	error  service.ErrorFunc
	tree   uci.Tree
	tm     *threadmanager.ThreadManager
	fs     filesystem.FileSystem
}

func (c *config) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	c.logger = logger
	c.error = e
	if err := servicehelpers.CheckMainDependencies(c, services); err != nil {
		return err
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{c.initSvcFileSystem}); err != nil {
		return err
	}
	c.tree = uci.NewTreeFromFs(afero.NewBasePathFs(nil, configPATH))
	c.tm = threadmanager.NewThreadManager(c.logger, c.error)
	c.tm.Load("commit", c.configCommitThread)
	c.tm.Start("commit")
	return nil
}

func (c *config) Stop() error {
	c.tm.StopAll()
	if err := c.commit(); err != nil {
		return err
	}
	return nil
}

func (c *config) Name() string { return "config" }

func (c *config) Dependencies() ([]string, []string) {
	return []string{"filesystem"}, []string{}
}

func (c *config) commit() error {
	return c.tree.Commit()
}

func (c *config) GetServiceConfig(s service.Service) *ServiceConfig {
	return newServiceConfig(c, s.Name())
}
