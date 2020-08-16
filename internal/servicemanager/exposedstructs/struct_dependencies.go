package exposedstructs

type Dependencies struct {
	Dependencies                     []string `json:"dependencies"`
	AddDependencies                  []string `json:"add_dependencies"`
	AddDependenciesSet               []string `json:"add_dependencies_set"`
	ReverseDependencies              []string `json:"rev_dependencies"`
	ReverseAdditionalDependencies    []string `json:"rev_add_dependencies"`
	ReverseAdditionalDependenciesSet []string `json:"rev_add_dependencies_set"`
}
