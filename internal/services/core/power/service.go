package power

import (
	"os/exec"
	"strconv"
	"time"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"github.com/gorilla/mux"
)

const (
	PERMISSION = "power"
)

type power struct {
	acm   acm.ACM
	web   web.Web
	mux   *mux.Router
	error service.ErrorFunc
}

func (p *power) Init(services map[string]service.Service, _ service.LoggerFunc, e service.ErrorFunc) {
	p.error = e
	if err := servicehelpers.CheckMainDependencies(p, services); err != nil {
		e(err)
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{"acm": p.initSvcAcm}); err != nil {
		e(err)
	}
}

func (p *power) Stop()        {}
func (p *power) Name() string { return "power" }
func (p *power) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{Deps: []string{"acm"}, ADeps: []string{"web"}}
}
func (p *power) SetAdditionalDependencies(services map[string]service.Service) {
	servicehelpers.InitializeAdditionalDependencies(services, servicehelpers.AdditionalServiceInitializers{"web": p.initSvcWeb})
}
func (p *power) UnsetAdditionalDependencies([]string) {}

func (p *power) Poweroff(token string) (error, bool) {
	if err, uae := p.acm.CheckTokenPermission(token, PERMISSION); err != nil {
		return err, uae
	}
	cmd := exec.Command("sudo", "poweroff")
	return cmd.Run(), false
}

func (p *power) Reboot(token string) (error, bool) {
	if err, uae := p.acm.CheckTokenPermission(token, PERMISSION); err != nil {
		return err, uae
	}
	cmd := exec.Command("sudo", "reboot")
	return cmd.Run(), false
}

func (p *power) SetWakeAlarm(t time.Time, token string) (error, bool) {
	if err, uae := p.acm.CheckTokenPermission(token, PERMISSION); err != nil {
		return err, uae
	}
	cmd := exec.Command("echo", ">", strconv.Itoa(int(t.Unix())))
	return cmd.Run(), false
}

func (p *power) UnsetWakeAlarm(token string) (error, bool) {
	return p.SetWakeAlarm(time.Unix(0, 0), token)
}
