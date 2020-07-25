package fwver

import (
	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/filehelpers"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
)

type fwver struct {
	rawver string
	semver semver.Version
	web    web.Web
	mux    *mux.Router
	fs     filesystem.FileSystem
}

func (f *fwver) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{f.initSvcFileSystem}); err != nil {
		return err
	}
	var err error
	f.rawver, err = filehelpers.ReadFirstLineString(f.fs, versionpath)
	if err != nil {
		return err
	}
	f.semver, err = semver.Make(f.rawver)
	if err != nil {
		return err
	}
	return nil
}

func (f *fwver) Name() string                       { return "fwver" }
func (f *fwver) Dependencies() ([]string, []string) { return []string{"filesystem"}, []string{"web"} }
func (f *fwver) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(f, services); err != nil {
		return err
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{f.initSvcWeb}); err != nil {
		return err
	}
	return nil
}

func (f *fwver) Get() semver.Version {
	return f.semver
}

func (f *fwver) GetString() string {
	return f.rawver
}
