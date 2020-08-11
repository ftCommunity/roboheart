package deviceinfo

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"net/http"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
)

func (d *deviceinfo) initSvcWeb(svc service.Service) {
	var ok bool
	d.web, ok = svc.(web.Web)
	if !ok {
		d.error(errors.New("Type assertion error"))
	}
	d.mux = d.web.RegisterServiceAPI(d)
	d.mux.HandleFunc("/platform", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, d.GetPlatform())
	})
	d.mux.HandleFunc("/device", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, d.GetDevice())
	})
}
