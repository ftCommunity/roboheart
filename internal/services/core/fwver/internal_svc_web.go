package fwver

import (
	"errors"
	"net/http"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
)

func (f *fwver) initSvcWeb(services servicehelpers.ServiceList) error {
	var ok bool
	f.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	f.mux = f.web.RegisterServiceAPI(f)
	f.mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, f.semver.String())
	})
	return nil
}
