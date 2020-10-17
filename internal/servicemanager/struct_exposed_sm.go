package servicemanager

//this must implement instance.ServiceManager
type exposed struct {
	sm *ServiceManager
}

func newExposed(sm *ServiceManager) *exposed {
	e := &exposed{}
	e.sm = sm
	return e
}
