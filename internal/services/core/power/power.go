package power

import (
	"errors"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
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

type Power interface {
	Poweroff(token string) error
	Reboot(token string) error
	SetWakeAlarm(t time.Time, token string) error
	UnsetWakeAlarm(token string) error
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

func (p *power) Stop() error                        { return nil }
func (p *power) Name() string                       { return "power" }
func (p *power) Dependencies() ([]string, []string) { return []string{"acm"}, []string{"web"} }
func (p *power) SetAdditionalDependencies(services map[string]service.Service) error {
	if err := servicehelpers.CheckAdditionalDependencies(p, services); err != nil {
		return err
	}
	var ok bool
	p.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	p.configureWeb()
	return nil
}
func (p *power) UnsetAdditionalDependencies() {}

func (p *power) configureWeb() {
	p.mux = p.web.RegisterServiceAPI(p)
	p.mux.HandleFunc("/poweroff", func(w http.ResponseWriter, r *http.Request) {
		data := &api.TokenRequest{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		if err := p.Poweroff(data.Token); err != nil {
			code := 500
			if acm.IsUserError(err) {
				code = 403
			}
			api.ErrorResponseWriter(w, code, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	}).Methods("POST")
	p.mux.HandleFunc("/reboot", func(w http.ResponseWriter, r *http.Request) {
		data := &api.TokenRequest{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		if err := p.Reboot(data.Token); err != nil {
			code := 500
			if acm.IsUserError(err) {
				code = 403
			}
			api.ErrorResponseWriter(w, code, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	}).Methods("POST")
	p.mux.HandleFunc("/wakealarm", func(w http.ResponseWriter, r *http.Request) {
		data := &wakeAlarmRequest{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		if err := p.SetWakeAlarm(time.Unix(data.Time, 0), data.Token); err != nil {
			code := 500
			if acm.IsUserError(err) {
				code = 403
			}
			api.ErrorResponseWriter(w, code, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	})
	p.mux.HandleFunc("/wakealarm", func(w http.ResponseWriter, r *http.Request) {
		data := &api.TokenRequest{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		if err := p.UnsetWakeAlarm(data.Token); err != nil {
			code := 500
			if acm.IsUserError(err) {
				code = 403
			}
			api.ErrorResponseWriter(w, code, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	}).Methods("DELETE")
}

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

func (p *power) SetWakeAlarm(t time.Time, token string) error {
	if err := acm.CheckTokenPermission(p.acm, token, PERMISSION); err != nil {
		return err
	}
	cmd := exec.Command("echo", ">", strconv.Itoa(int(t.Unix())))
	return cmd.Run()
}

func (p *power) UnsetWakeAlarm(token string) error {
	return p.SetWakeAlarm(time.Unix(0, 0), token)
}

var Service = new(power)
