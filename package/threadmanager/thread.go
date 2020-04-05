package threadmanager

import (
	"errors"
	"sync"

	"github.com/ftCommunity/roboheart/internal/service"
)

type threadfunc func(service.LoggerFunc, service.ErrorFunc, chan interface{}, chan interface{})

type thread struct {
	stopc, stopped chan interface{}
	logger         service.LoggerFunc
	error          service.ErrorFunc
	state          bool
	lock           sync.Mutex
	f              threadfunc
}

func (t *thread) start() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.state {
		return errors.New("Thread already running")
	}
	t.call()
	t.state = true
	return nil
}

func (t *thread) stop() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.state {
		t.stopc <- struct{}{}
		<-t.stopped
		t.logger("Stopped")
		t.state = false
	}
}

func (t *thread) restart() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.state {
		t.stopc <- struct{}{}
		<-t.stopped
		t.logger("Stopped")
	}
	t.call()
	t.state = true
}

func (t *thread) call() {
	go t.f(t.logger, t.error, t.stopc, t.stopped)
	t.logger("Started")
}

func newThread(f threadfunc, logger service.LoggerFunc, e service.ErrorFunc) *thread {
	t := new(thread)
	t.f = f
	t.logger = logger
	t.error = e
	t.stopc = make(chan interface{})
	t.stopped = make(chan interface{})
	return t
}
