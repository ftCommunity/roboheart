package deviceinfo

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/filehelpers"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
)

type deviceinfo struct {
	web              web.Web
	mux              *mux.Router
	fs               filesystem.FileSystem
	platform, device string
	error            service.ErrorFunc
}

func (d *deviceinfo) Init(services map[string]service.Service, _ service.LoggerFunc, e service.ErrorFunc) {
	d.error = e
	if err := servicehelpers.CheckMainDependencies(d, services); err != nil {
		e(err)
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{"filesystem": d.initSvcFileSystem}); err != nil {
		e(err)
	}
	var err error
	d.platform, err = filehelpers.ReadFirstLineString(d.fs, platformpath)
	if err != nil {
		e(err)
	}
	d.device, err = filehelpers.ReadFirstLineString(d.fs, devicepath)
	if err != nil {
		e(err)
	}
}

func (d *deviceinfo) Stop() {}

func (d *deviceinfo) Name() string { return "deviceinfo" }
func (d *deviceinfo) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{Deps: []string{"filesystem"}, ADeps: []string{"web"}}
}
func (d *deviceinfo) SetAdditionalDependencies(services map[string]service.Service) {
	servicehelpers.InitializeAdditionalDependencies(services, servicehelpers.AdditionalServiceInitializers{"web": d.initSvcWeb})
}

func (d *deviceinfo) UnsetAdditionalDependencies([]string) {}

func (d *deviceinfo) GetPlatform() string { return d.platform }
func (d *deviceinfo) GetDevice() string   { return d.device }
