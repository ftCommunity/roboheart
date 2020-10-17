package servicemanager

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
	for {
		select {
		case <-sm.workercheck:
			{
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
	}
}

func (sm *ServiceManager) workertask() {
	sm.serviceslock.Lock()
	defer sm.serviceslock.Unlock()
	for _, ss := range sm.services {
		for _, is := range ss.instances {
			if !is.running {
				is.instance.base.Init(is.logger, is.error)
			}
			if di := is.instance.depending; di != nil {
				ndeps := di.Dependencies()
				o, n, _ := is.deps.deps.Compare(ndeps)
				for _, od := range o {
					di.UnsetDependency(od)
					is.deps.deps.Delete(od)
					sm.services[od.Name].get(od).deps.rdeps.Delete(is.id)
				}
			newdeps:
				for _, nd := range n {
					if _, ok := sm.services[nd.Name]; !ok {
						continue newdeps
					}
					if ni := sm.services[nd.Name].get(nd); ni == nil {
						if sm.newInstance(nd) != nil {
							continue newdeps
						}
						ni = sm.services[nd.Name].get(nd)
					} else {
						is.deps.deps.Add(nd)
						ni.deps.rdeps.Add(is.id)
						di.SetDependency(ni.instance.base)
					}
				}
			}
		}
	}
}
