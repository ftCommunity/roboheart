package vncserver

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/labstack/echo/v4"
)

func (v *vncserver) initSvcWeb(svc service.Service) {
	var ok bool
	v.web, ok = svc.(web.Web)
	if !ok {
		v.error(errors.New("Type assertion error"))
	}
	v.web.RegisterServiceAPI(v)
}

func (v *vncserver) deinitSvcWeb() {
	v.web.UnregisterServiceAPI(v)
}

func (v *vncserver) WebRegisterRoutes(group *echo.Group) {
	group.GET("/state", func(c echo.Context) error {
		api.ResponseWriter(c, state{v.state})
		return nil
	})
	group.POST("/state", func(c echo.Context) error {
		data := &stateSet{}
		if !api.RequestLoader(c, data) {
			return nil
		}
		var f func(string) (error, bool)
		if data.State {
			f = v.StartVNC
		} else {
			f = v.StopVNC
		}
		if err, _ := f(data.Token); err != nil {
			api.ErrorResponseWriter(c, 500, err)
		} else {
			api.ResponseWriter(c, nil)
		}
		return nil
	})
	group.GET("/autostart", func(c echo.Context) error {
		api.ResponseWriter(c, autostartState{v.GetAutostart()})
		return nil
	})
	group.POST("/autostart", func(c echo.Context) error {
		data := &autostartSet{}
		if !api.RequestLoader(c, data) {
			return nil
		}
		if err, _ := v.SetAutostart(data.Token, data.Autostart); err != nil {
			api.ErrorResponseWriter(c, 500, err)
		} else {
			api.ResponseWriter(c, nil)
		}
		return nil
	})
}
