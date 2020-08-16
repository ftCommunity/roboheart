package exposedstructs

type ServiceInfo struct {
	Name         string       `json:"name"`
	Running      bool         `json:"runing"`
	Dependencies Dependencies `json:"dependencies"`
	Builtin      bool         `json:"builtin"`
}
