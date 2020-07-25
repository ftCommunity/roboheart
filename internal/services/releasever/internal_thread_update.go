package relver

import (
	"time"

	"github.com/ftCommunity/roboheart/internal/service"
)

func (r *relver) updateThread(logger service.LoggerFunc, e service.ErrorFunc, stop, stopped chan interface{}) {
	if err := r.getReleaseData(); err != nil {
		logger(err)
	}
	for {
		select {
		case <-stop:
			{
				stopped <- struct{}{}
				return
			}
		case <-time.After(15 * time.Minute):
			{
				if err := r.getReleaseData(); err != nil {
					logger(err)
				}
			}
		}
	}
}
