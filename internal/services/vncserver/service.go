package vncserver

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/procrunner"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
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

func (v *vncserver) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	v.logger = logger
	v.error = e
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{v.initSvcAcm, v.initSvcConfig}); err != nil {
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
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{v.initSvcWeb}); err != nil {
		return err
	}
	return nil
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
