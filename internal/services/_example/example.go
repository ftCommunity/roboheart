package example

import (
	"time"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/threadmanager"
)

type example struct {
	logger service.LoggerFunc
	error  service.ErrorFunc
	tm     *threadmanager.ThreadManager
}

type Example interface {
	MethodOne()
	MethodTwo(_ string, _ int) bool
}

func (ex *example) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	ex.logger = logger
	ex.error = e
	ex.tm = threadmanager.NewThreadManager(ex.logger, ex.error)
	ex.tm.Load("cleanup", ex.exampleThread)
	ex.tm.Start("cleanup")
	return nil
}

func (ex *example) Stop() error                                                { return nil }
func (ex *example) Name() string                                               { return "example" }
func (ex *example) Dependencies() ([]string, []string)                         { return []string{}, []string{} }
func (ex *example) SetAdditionalDependencies(map[string]service.Service) error { return nil }
func (ex *example) UnsetAdditionalDependencies()                               {}

func (ex *example) exampleThread(logger service.LoggerFunc, e service.ErrorFunc, stop, stopped chan interface{}) {
	//for a normal "do every x seconds"-thread you should not need to change too much
	//if your task should run once beforehand you might just put it right before the for loop
	for {
		select {
		case <-stop:
			{
				stopped <- struct{}{}
				return
			}
		//change time here
		case <-time.After(5 * time.Second):
			{
				//task
			}
		}
	}
}

func (ex *example) MethodOne()                     {}
func (ex *example) MethodTwo(_ string, _ int) bool { return true }

var Service = new(example)
