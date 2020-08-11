package power

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"net/http"
	"time"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
)

func (p *power) initSvcWeb(svc service.Service) {
	var ok bool
	p.web, ok = svc.(web.Web)
	if !ok {
		p.error(errors.New("Type assertion error"))
	}
	p.mux = p.web.RegisterServiceAPI(p)
	p.mux.HandleFunc("/poweroff", func(w http.ResponseWriter, r *http.Request) {
		data := &api.TokenRequest{}
		if !api.RequestLoader(r, w, data) {
			return
		}
		if err, uae := p.Poweroff(data.Token); err != nil {
			code := 500
			if uae {
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
		if err, uae := p.Reboot(data.Token); err != nil {
			code := 500
			if uae {
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
		if err, uae := p.SetWakeAlarm(time.Unix(data.Time, 0), data.Token); err != nil {
			code := 500
			if uae {
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
		if err, uae := p.UnsetWakeAlarm(data.Token); err != nil {
			code := 500
			if uae {
				code = 403
			}
			api.ErrorResponseWriter(w, code, err)
		} else {
			api.ResponseWriter(w, nil)
		}
	}).Methods("DELETE")
}
