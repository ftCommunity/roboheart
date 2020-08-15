package fwver

import (
	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/filehelpers"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

type fwver struct {
	rawver string
	semver semver.Version
	web    web.Web
	fs     filesystem.FileSystem
	error  service.ErrorFunc
}

func (f *fwver) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) {
	f.error = e
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{"filesystem": f.initSvcFileSystem}); err != nil {
		e(err)
	}
	var err error
	f.rawver, err = filehelpers.ReadFirstLineString(f.fs, versionpath)
	if err != nil {
		e(err)
	}
	f.semver, err = semver.Make(f.rawver)
	if err != nil {
		e(err)
	}
}

func (f *fwver) Stop()        {}
func (f *fwver) Name() string { return "fwver" }
func (f *fwver) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{Deps: []string{"filesystem"}, ADeps: []string{"web"}}
}
func (f *fwver) SetAdditionalDependencies(services map[string]service.Service) {
	servicehelpers.InitializeAdditionalDependencies(services, servicehelpers.AdditionalServiceInitializers{"web": f.initSvcWeb})
}
func (f *fwver) UnsetAdditionalDependencies(services []string) {
	servicehelpers.DeinitAdditionalDependencies(services, servicehelpers.AdditionalServiceDeinitializers{"web": f.deinitSvcWeb})
}

func (f *fwver) Get() semver.Version {
	return f.semver
}

func (f *fwver) GetString() string {
	return f.rawver
}
