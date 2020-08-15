package fwver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/labstack/echo/v4"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
)

func (f *fwver) initSvcWeb(svc service.Service) {
	var ok bool
	f.web, ok = svc.(web.Web)
	if !ok {
		f.error(errors.New("Type assertion error"))
	}
	f.web.RegisterServiceAPI(f)
}

func (f *fwver) deinitSvcWeb() {
	f.web.UnregisterServiceAPI(f)
}

func (f *fwver) WebRegisterRoutes(group *echo.Group) {
	group.GET("/version", func(c echo.Context) error {
		api.ResponseWriter(c, f.semver.String())
		return nil
	})
}
