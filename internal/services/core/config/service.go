package config

import (
	"github.com/digineo/go-uci"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/threadmanager"
)

type config struct {
	logger service.LoggerFunc
	error  service.ErrorFunc
	tree   uci.Tree
	tm     *threadmanager.ThreadManager
}

func (c *config) Init(_ map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	c.logger = logger
	c.error = e
	c.tree = uci.NewTree(configPATH)
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

func (c *config) commit() error {
	return c.tree.Commit()
}

func (c *config) GetServiceConfig(s service.Service) *ServiceConfig {
	return newServiceConfig(c, s.Name())
}
