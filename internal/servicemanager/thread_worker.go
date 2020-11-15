package servicemanager

import (
	"time"
)

func (sm *ServiceManager) triggerWorker() {
	sm.workercalllock.Lock()
	defer sm.workercalllock.Unlock()
	sm.workercall = true
	sm.workercheck <- struct{}{}
}

func (sm *ServiceManager) stopWorker() {
	sm.workerstop = true
	sm.workercheck <- struct{}{}
}

func (sm *ServiceManager) worker() {
	for range sm.workercheck {
		if sm.workerstop {
			sm.workerstop = false
			sm.wg.Done()
			return
		}
		sm.workercalllock.Lock()
		if sm.workercall {
			sm.workercalllock.Unlock()
			sm.workertask()
		} else {
			sm.workercalllock.Unlock()
		}
	}
}

func (sm *ServiceManager) workertask() {
workerloop:
	for {
		sm.workerrunninglock.Lock()
		sm.workerrunning = true
		sm.workerrunninglock.Unlock()
		sm.serviceslock.Lock()
		for _, ss := range sm.services {
			for _, is := range ss.instances {
				if !is.running {
					startok := make(chan interface{})
					go func() {
						is.start()
						startok <- struct{}{}
					}()
					select {
					case <-startok:
						{
						}
					case <-sm.workerabort:
						{
							sm.workerrunninglock.Lock()
							sm.workerrunning = false
							sm.workerrunninglock.Unlock()
							sm.serviceslock.Unlock()
							continue workerloop
						}
					}
					is.running = true
				}
				is.updateDependencies()
				if !is.startup {
					if is.created.Add(INSTANCE_NO_REASON_TIMEOUT).Before(time.Now()) &&
						is.lastrdep.Add(INSTANCE_NO_REASON_TIMEOUT).Before(time.Now()) &&
						len(*is.deps.rdeps) == 0 {
						for _, di := range *is.deps.deps {
							is.instance.depending.UnsetDependency(di)
							sm.get(di).deps.rdeps.Delete(is.id)
						}
						stopok := make(chan interface{})
						go func() {
							is.stop()
							stopok <- struct{}{}
						}()
						select {
						case <-stopok:
							{
								delete(ss.instances, is.id.Instance)
							}
						case <-sm.workerabort:
							{
								delete(ss.instances, is.id.Instance)
								sm.workerrunninglock.Lock()
								sm.workerrunning = false
								sm.workerrunninglock.Unlock()
								sm.serviceslock.Unlock()
								continue workerloop
							}
						}
					}
				}
			}
		}
		sm.workerrunninglock.Lock()
		sm.workerrunning = false
		sm.workerrunninglock.Unlock()
		sm.serviceslock.Unlock()
		return
	}
}
