package deviceinfo

import (
	"errors"
	"net/http"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/servicehelpers"

	"github.com/ftCommunity/roboheart/package/api"
)

func (d *deviceinfo) initSvcWeb(services servicehelpers.ServiceList) error {
	var ok bool
	d.web, ok = services["web"].(web.Web)
	if !ok {
		return errors.New("Type assertion error")
	}
	d.mux = d.web.RegisterServiceAPI(d)
	d.mux.HandleFunc("/platform", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, d.GetPlatform())
	})
	d.mux.HandleFunc("/device", func(w http.ResponseWriter, _ *http.Request) {
		api.ResponseWriter(w, d.GetDevice())
	})
	return nil
}
