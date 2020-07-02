package pkgmanager

import "github.com/ftCommunity/roboheart/package/packages"

type pkg struct {
	packages.Package
	Dependencies map[string]dependency `json:"dependencies"`
}

type dependency struct {
	packages.Dependency
	satisfied bool
}
