package pkgmanager

import "github.com/ftCommunity/roboheart/package/packages"

type pkg struct {
	packages.Package
	Dependencies map[string]dependency `json:"dependencies"`
}

type dependency struct {
	packages.Dependency
	OneOf     []dependencyOption `json:"oneof"`
	satisfied bool
}

type dependencyOption struct {
	packages.DependencyOption
	used bool
}
