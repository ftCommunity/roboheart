package robotime

import "github.com/ftCommunity/roboheart/internal/service"
import "time"

type robotime struct {
}

func (r robotime) Init(map[string]service.Service, service.LoggerFunc, service.ErrorFunc) {}

func (r robotime) Name() string {
	return "robotime"
}

func (r robotime) Stop() {}

func (r robotime) Now() time.Time {
	return time.Now()
}

func (r robotime) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
