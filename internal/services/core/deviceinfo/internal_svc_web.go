package deviceinfo

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/labstack/echo/v4"
)

func (d *deviceinfo) initSvcWeb(svc service.Service) {
	var ok bool
	d.web, ok = svc.(web.Web)
	if !ok {
		d.error(errors.New("Type assertion error"))
	}
	d.web.RegisterServiceAPI(d)
}

func (d *deviceinfo) deinitSvcWeb() {
	d.web.UnregisterServiceAPI(d)
}

func (d *deviceinfo) WebRegisterRoutes(group *echo.Group) {
	group.GET("/platform", func(c echo.Context) error {
		api.ResponseWriter(c, d.GetPlatform())
		return nil
	})
	group.GET("/device", func(c echo.Context) error {
		api.ResponseWriter(c, d.GetDevice())
		return nil
	})
}
