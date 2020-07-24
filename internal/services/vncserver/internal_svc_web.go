package vncserver

import (
	"errors"
	"net/http"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (v *vncserver) initSvcWeb(services servicehelpers.ServiceList) error {
	var ok bool
	v.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
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
	return nil
}
