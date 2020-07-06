package pkgmanager

import (
	"github.com/ftCommunity/roboheart/package/packages"
	"sort"
)

const (
	STATE_TODO       = 0
	STATE_DONE       = 1
	STATE_IMPOSSIBLE = 2
)

type extendedPackage struct {
	packages.Package
	Variants map[string]extendedVariant
}

type extendedVariant struct {
	packages.Variant
	Dependencies map[string]dependency `json:"dependencies"`
	state        int
}

type extendedVariantList []extendedVariant

func (e extendedVariantList) Len() int {
	return len(e)
}

func (e extendedVariantList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e extendedVariantList) Less(i, j int) bool {
	return e[i].Version.LT(*e[j].Version.Version)
}

func (e *extendedVariantList) Sort() {
	sort.Sort(e)
}

type dependency struct {
	packages.Dependency
	OneOf     []dependencyOption `json:"oneof"`
	satisfied bool
}

type dependencyOption struct {
	packages.DependencyOption
	Wants   *dependencyOptionWants   `json:"wants"`
	Package *dependencyOptionPackage `json:"package"`
}

type dependencyOptionExtension struct {
	used         bool
	pkg, variant string
}

type dependencyOptionWants struct {
	packages.DependencyOptionWants
	dependencyOptionExtension
}

type dependencyOptionPackage struct {
	packages.DependencyOptionPackage
	dependencyOptionExtension
}
