package config

import (
	"time"

	"github.com/digineo/go-uci"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/threadmanager"
)

const (
	configPATH = "/etc/config"
)

type config struct {
	logger service.LoggerFunc
	error  service.ErrorFunc
	tree   uci.Tree
	tm     *threadmanager.ThreadManager
}

type Config interface {
	GetTree() uci.Tree
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

func (c *config) Dependencies() ([]string, []string) { return []string{}, []string{} }

func (c *config) SetAdditionalDependencies(map[string]service.Service) error { return nil }

func (c *config) commit() error {
	return c.tree.Commit()
}

func (c *config) configCommitThread(logger service.LoggerFunc, e service.ErrorFunc, stop, stopped chan interface{}) {
	for {
		select {
		case <-stop:
			{
				stopped <- struct{}{}
				return
			}
		case <-time.After(5 * time.Second):
			{
				if err := c.commit(); err != nil {
					e(err)
				}
			}
		}
	}
}

func (c *config) UnsetAdditionalDependencies() {}

func (c *config) GetTree() uci.Tree { return c.tree }

var Service = new(config)
