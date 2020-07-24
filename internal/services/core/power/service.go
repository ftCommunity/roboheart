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
	acm acm.ACM
	web web.Web
	mux *mux.Router
}

func (p *power) Init(services map[string]service.Service, _ service.LoggerFunc, _ service.ErrorFunc) error {
	if err := servicehelpers.CheckMainDependencies(p, services); err != nil {
		return err
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{p.initSvcAcm}); err != nil {
		return err
	}
	return nil
}

func (p *power) Stop() error                        { return nil }
func (p *power) Name() string                       { return "power" }
func (p *power) Dependencies() ([]string, []string) { return []string{"acm"}, []string{"web"} }
func (p *power) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(p, services); err != nil {
		return err
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{p.initSvcWeb}); err != nil {
		return err
	}
	return nil
}
func (p *power) UnsetAdditionalDependencies() {}

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
