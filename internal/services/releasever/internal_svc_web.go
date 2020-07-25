package relver

import (
	"errors"
	"net/http"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (r *relver) initSvcWeb(services servicehelpers.ServiceList) error {
	var ok bool
	r.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	r.mux = r.web.RegisterServiceAPI(r)
	r.mux.HandleFunc("/release", func(w http.ResponseWriter, _ *http.Request) {
		defer r.lock.Unlock()
		r.lock.Lock()
		if r.release != nil {
			api.ResponseWriter(w, r.release)
		} else {
			api.ErrorResponseWriter(w, 503, errors.New("Version information not available"))
		}
	})
	r.mux.HandleFunc("/prerelease", func(w http.ResponseWriter, _ *http.Request) {
		defer r.lock.Unlock()
		r.lock.Lock()
		if r.prerelease != nil {
			api.ResponseWriter(w, r.prerelease)
		} else {
			api.ErrorResponseWriter(w, 503, errors.New("Version information not available"))
		}
	})
	r.mux.HandleFunc("/releases", func(w http.ResponseWriter, _ *http.Request) {
		defer r.lock.Unlock()
		r.lock.Lock()
		api.ResponseWriter(w, r.releases)
	})
	return nil
}
