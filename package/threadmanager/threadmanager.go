package threadmanager

import (
	"errors"
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
	"sync"
)

type ThreadManager struct {
	threads  map[string]*thread
	lock     sync.Mutex
	lockdown bool
	logger   instance.LoggerFunc
	error    instance.ErrorFunc
}

func (tm *ThreadManager) Load(id string, f threadfunc) error {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if tm.lockdown {
		return errors.New("Threadmanager stopped")
	}
	if _, set := tm.threads[id]; set {
		return errors.New("Thread ID already registered")
	}
	tm.threads[id] = newThread(f, tm.genThreadLogger(id), tm.genThreadError(id))
	return nil
}

func (tm *ThreadManager) getThread(id string) (*thread, error) {
	if t, ok := tm.threads[id]; ok {
		return t, nil
	} else {
		return nil, errors.New("Thread not found")
	}
}

func (tm *ThreadManager) Start(id string) error {
	t, err := tm.getThread(id)
	if err != nil {
		return err
	}
	return t.start()
}

func (tm *ThreadManager) Stop(id string) error {
	t, err := tm.getThread(id)
	if err != nil {
		return err
	}
	t.stop()
	return nil
}

func (tm *ThreadManager) Restart(id string) error {
	t, err := tm.getThread(id)
	if err != nil {
		return err
	}
	t.restart()
	return nil
}

func (tm *ThreadManager) StopAll() {
	tm.lock.Lock()
	tm.lockdown = true
	tm.lock.Unlock()
	for _, t := range tm.threads {
		t.stop()
	}
}

func (tm *ThreadManager) genThreadLogger(tn string) instance.LoggerFunc {
	return func(v ...interface{}) {
		tm.logger(append([]interface{}{"Thread:", tn + ":"}, v...)...)
	}
}

func (tm *ThreadManager) genThreadError(tn string) instance.ErrorFunc {
	return func(v ...interface{}) {
		tm.error(append([]interface{}{"Thread:", tn + ":"}, v...)...)
	}
}

func NewThreadManager(logger instance.LoggerFunc, e instance.ErrorFunc) *ThreadManager {
	tm := new(ThreadManager)
	tm.threads = make(map[string]*thread)
	tm.logger = logger
	tm.error = e
	return tm
}
