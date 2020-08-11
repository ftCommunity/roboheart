package fwver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"net/http"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
)

func (f *fwver) initSvcWeb(svc service.Service) {
	var ok bool
	f.web, ok = svc.(web.Web)
	if !ok {
		f.error(errors.New("Type assertion error"))
	}
	f.mux = f.web.RegisterServiceAPI(f)
	f.mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, f.semver.String())
	})
}
