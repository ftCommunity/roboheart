package servicemanager

import "time"

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
	sm.serviceslock.Lock()
	defer sm.serviceslock.Unlock()
	for _, ss := range sm.services {
		for _, is := range ss.instances {
			if !is.running {
				is.instance.base.Start()
				is.running = true
			}
			is.updateDependencies()
			if !is.startup {
				if is.lastrdep.Add(SERVICE_NO_REASON_TIMEOUT).Before(time.Now()) && len(*is.deps.rdeps) == 0 {
					for _, di := range *is.deps.deps {
						is.instance.depending.UnsetDependency(di)
						sm.get(di).deps.rdeps.Delete(is.id)
					}
					is.getBase().Stop()
					delete(ss.instances, is.id.Instance)
				}
			}
		}
	}
}
