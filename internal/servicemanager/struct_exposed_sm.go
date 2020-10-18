package servicemanager

//this must implement instance.ServiceManager
type exposed struct {
	sm *ServiceManager
}

func (e *exposed) LoadFromPlugin(path string) error {
	return e.sm.loadFromPlugin(path)
}

func newExposed(sm *ServiceManager) *exposed {
	e := &exposed{}
	e.sm = sm
	return e
}
