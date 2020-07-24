package fwver

import (
	"io/ioutil"
	"strings"

	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
)

type fwver struct {
	rawver string
	semver semver.Version
	web    web.Web
	mux    *mux.Router
}

func (f *fwver) Init(map[string]service.Service, service.LoggerFunc, service.ErrorFunc) error {
	raw, err := ioutil.ReadFile("/etc/fw-ver.txt")
	if err != nil {
		return err
	}
	f.rawver = strings.Split(string(raw), "\n")[0]
	f.semver, err = semver.Make(f.rawver)
	if err != nil {
		return err
	}
	return nil
}

func (f *fwver) Name() string                       { return "fwver" }
func (f *fwver) Dependencies() ([]string, []string) { return []string{}, []string{"web"} }
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
