package config

import (
	"time"

	"github.com/ftCommunity/roboheart/internal/service"
)

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
