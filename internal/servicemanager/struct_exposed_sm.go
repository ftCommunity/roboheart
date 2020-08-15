package servicemanager

//this must implement service.ServiceManager
type exposed struct {
	sm *ServiceManager
}

func newExposed(sm *ServiceManager) *exposed {
	e := &exposed{}
	e.sm = sm
	return e
}
