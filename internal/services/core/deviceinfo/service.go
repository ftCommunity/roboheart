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
}

func (d *deviceinfo) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	if err := servicehelpers.CheckMainDependencies(d, services); err != nil {
		return err
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{d.initSvcFileSystem}); err != nil {
		return err
	}
	var err error
	d.platform, err = filehelpers.ReadFirstLineString(d.fs, platformpath)
	if err != nil {
		return err
	}
	d.device, err = filehelpers.ReadFirstLineString(d.fs, devicepath)
	if err != nil {
		return err
	}
	return nil
}

func (d *deviceinfo) Name() string { return "deviceinfo" }
func (d *deviceinfo) Dependencies() ([]string, []string) {
	return []string{"filesystem"}, []string{"web"}
}
func (d *deviceinfo) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(d, services); err != nil {
		return err
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{d.initSvcWeb}); err != nil {
		return err
	}
	return nil
}

func (d *deviceinfo) GetPlatform() string { return d.platform }
func (d *deviceinfo) GetDevice() string   { return d.device }
