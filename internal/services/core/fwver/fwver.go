package fwver

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
)

type fwver struct {
	rawver string
	semver semver.Version
	web    web.Web
	mux    *mux.Router
}

type FWVer interface {
	Get() semver.Version
	GetString() string
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

func (f *fwver) Stop() error                        { return nil }
func (f *fwver) Name() string                       { return "fwver" }
func (f *fwver) Dependencies() ([]string, []string) { return []string{}, []string{"web"} }
func (f *fwver) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(f, services); err != nil {
		return err
	}
	var ok bool
	f.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	f.configureWeb()
	return nil
}

func (f *fwver) UnsetAdditionalDependencies() {}

func (f *fwver) configureWeb() {
	f.mux = f.web.RegisterServiceAPI(f)
	f.mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, f.semver.String())
	})
}

func (f *fwver) Get() semver.Version {
	return f.semver
}

func (f *fwver) GetString() string {
	return f.rawver
}

var Service = new(fwver)
