package vncserver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/ftCommunity/roboheart/package/procrunner"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	PERMISSION              = "vncserver"
	CONFIG_SECTION          = ""
	CONFIG_TYPE             = "vncserver"
	CONFIG_AUTOSTART_OPTION = "autostart"
)

type vncserver struct {
	logger  service.LoggerFunc
	error   service.ErrorFunc
	proc    *procrunner.ProcRunner
	acm     acm.ACM
	web     web.Web
	mux     *mux.Router
	config  config.Config
	sconfig *config.ServiceConfig
	state   bool
}

type VNCServer interface {
	Start(token string) (error, bool)
	Stop(token string) (error, bool)
	GetAutostart() bool
	SetAutostart(token string, autostart bool) (error, bool)
}

func (v *vncserver) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	v.logger = logger
	v.error = e
	var ok bool
	v.acm, ok = services["acm"].(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	v.acm.RegisterPermission(PERMISSION, map[string]bool{"user": true, "app": false}, map[string]string{})
	v.config, ok = services["config"].(config.Config)
	if !ok {
		return errors.New("Type assertion error")
	}
	v.sconfig = v.config.GetServiceConfig(v)
	if err := v.sconfig.AddSection(CONFIG_SECTION, CONFIG_TYPE); err != nil {
		return err
	}
	v.proc = procrunner.NewProcRunner("framebuffer-vncserver", "-f", "/dev/fb0", "-t", "/dev/input/event0")
	if v.GetAutostart() {
		v.start()
	}
	return nil
}

func (v *vncserver) Name() string { return "vncserver" }

func (v *vncserver) Dependencies() ([]string, []string) {
	return []string{"acm", "config"}, []string{"web"}
}

func (v *vncserver) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(v, services); err != nil {
		return err
	}
	var ok bool
	v.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	v.configureWeb()
	return nil
}

func (v *vncserver) configureWeb() {
	v.mux = v.web.RegisterServiceAPI(v)
	v.mux.HandleFunc("/state", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, state{v.state})
	}).Methods("GET")
	v.mux.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		data := &stateSet{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		var f func(string) (error, bool)
		if data.State {
			f = v.Start
		} else {
			f = v.Stop
		}
		if err, _ := f(data.Token); err != nil {
			api.ErrorResponseWriter(w, 500, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	}).Methods("POST")
	v.mux.HandleFunc("/autostart", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, autostartState{v.GetAutostart()})
	}).Methods("GET")
	v.mux.HandleFunc("/autostart", func(w http.ResponseWriter, r *http.Request) {
		data := &autostartSet{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		if err, _ := v.SetAutostart(data.Token, data.Autostart); err != nil {
			api.ErrorResponseWriter(w, 500, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	}).Methods("POST")
}

func (v *vncserver) onCrash(c int) {
	v.logger("Crashed with code", c)
	v.logger("Restarting process")
}

func (v *vncserver) Start(token string) (error, bool) {
	if err, uae := v.acm.CheckTokenPermission(token, PERMISSION); err != nil {
		return err, uae
	}
	return v.start(), false
}

func (v *vncserver) Stop(token string) (error, bool) {
	if err, uae := v.acm.CheckTokenPermission(token, PERMISSION); err != nil {
		return err, uae
	}
	if err := v.proc.Stop(); err != nil {
		return err, false
	}
	return nil, false
}

func (v *vncserver) GetAutostart() bool {
	return v.sconfig.GetBoolDefault(CONFIG_SECTION, CONFIG_AUTOSTART_OPTION, true)
}

func (v *vncserver) SetAutostart(token string, autostart bool) (error, bool) {
	if err, uae := v.acm.CheckTokenPermission(token, PERMISSION); err != nil {
		return err, uae
	}
	v.sconfig.Set(CONFIG_SECTION, CONFIG_AUTOSTART_OPTION, "1")
	return nil, false
}

func (v *vncserver) start() error {
	v.proc.SetOnAutoRestartCallback(v.onCrash)
	if err := v.proc.Start(); err != nil {
		return err
	}
	return nil
}

var Service = new(vncserver)
