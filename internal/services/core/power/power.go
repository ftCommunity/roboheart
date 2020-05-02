package power

import (
	"errors"
	"os/exec"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

const (
	PERMISSION = "power"
)

type power struct {
	acm acm.ACM
}

type Power interface {
	Poweroff(token string) error
	Reboot(token string) error
}

func (p *power) Init(services map[string]service.Service, _ service.LoggerFunc, _ service.ErrorFunc) error {
	if err := servicehelpers.CheckMainDependencies(p, services); err != nil {
		return err
	}
	sacm, ok := services["acm"].(acm.ACM)
	if !ok {
		return errors.New("Type assertion error")
	}
	p.acm = sacm
	p.acm.RegisterPermission(PERMISSION, map[string]bool{"root": true, "user": true, "app": false})
	return nil
}

func (p *power) Stop() error                                                { return nil }
func (p *power) Name() string                                               { return "power" }
func (p *power) Dependencies() ([]string, []string)                         { return []string{"acm"}, []string{} }
func (p *power) SetAdditionalDependencies(map[string]service.Service) error { return nil }
func (p *power) UnsetAdditionalDependencies(s chan interface{})             { s <- struct{}{} }

func (p *power) Poweroff(token string) error {
	if err := acm.CheckTokenPermission(p.acm, token, PERMISSION); err != nil {
		return err
	}
	cmd := exec.Command("sudo", "poweroff")
	return cmd.Run()
}

func (p *power) Reboot(token string) error {
	if err := acm.CheckTokenPermission(p.acm, token, PERMISSION); err != nil {
		return err
	}
	cmd := exec.Command("sudo", "reboot")
	return cmd.Run()
}

var Service = new(power)
