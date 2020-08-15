package relver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/labstack/echo/v4"
)

func (r *relver) initSvcWeb(svc service.Service) {
	var ok bool
	r.web, ok = svc.(web.Web)
	if !ok {
		r.error(errors.New("Type assertion error"))
	}
	r.web.RegisterServiceAPI(r)
}

func (r *relver) deinitSvcWeb() {
	r.web.UnregisterServiceAPI(r)
}

func (r *relver) WebRegisterRoutes(group *echo.Group) {
	group.GET("/release", func(c echo.Context) error {
		defer r.lock.Unlock()
		r.lock.Lock()
		if r.release != nil {
			api.ResponseWriter(c, r.release)
		} else {
			api.ErrorResponseWriter(c, 503, errors.New("Version information not available"))
		}
		return nil
	})
	group.GET("/prerelease", func(c echo.Context) error {
		defer r.lock.Unlock()
		r.lock.Lock()
		if r.prerelease != nil {
			api.ResponseWriter(c, r.prerelease)
		} else {
			api.ErrorResponseWriter(c, 503, errors.New("Version information not available"))
		}
		return nil
	})
	group.GET("/releases", func(c echo.Context) error {
		defer r.lock.Unlock()
		r.lock.Lock()
		api.ResponseWriter(c, r.releases)
		return nil
	})
}
