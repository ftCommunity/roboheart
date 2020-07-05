package deviceinfo

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
	"net/http"
)

type deviceinfo struct {
	web web.Web
	mux *mux.Router
}

type DeviceInfo interface {
	GetPlatform() string
	GetDevice() string
}

func (d *deviceinfo) Init(map[string]service.Service, service.LoggerFunc, service.ErrorFunc) error {
	return nil
}

func (d *deviceinfo) Name() string                       { return "deviceinfo" }
func (d *deviceinfo) Dependencies() ([]string, []string) { return []string{}, []string{"web"} }
func (d *deviceinfo) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(d, services); err != nil {
		return err
	}
	var ok bool
	d.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	d.configureWeb()
	return nil
}

func (d *deviceinfo) configureWeb() {
	d.mux = d.web.RegisterServiceAPI(d)
	d.mux.HandleFunc("/platform", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, d.GetPlatform())
	})
	d.mux.HandleFunc("/device", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, d.GetDevice())
	})
}

func (d *deviceinfo) GetPlatform() string { return "arm" }
func (d *deviceinfo) GetDevice() string   { return "ft-txt" }

var Service = new(deviceinfo)
