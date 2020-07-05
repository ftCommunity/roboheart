package pkgmanager

import (
	"github.com/ftCommunity/roboheart/package/packages"
)

type extendedPackage struct {
	packages.Package
	Variants map[string]extendedVariant
}

type extendedVariant struct {
	packages.Variant
	Dependencies map[string]dependency `json:"dependencies"`
	done         bool
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
