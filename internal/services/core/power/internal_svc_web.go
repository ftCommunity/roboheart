package power

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/labstack/echo/v4"
	"time"

	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
)

func (p *power) initSvcWeb(svc service.Service) {
	var ok bool
	p.web, ok = svc.(web.Web)
	if !ok {
		p.error(errors.New("Type assertion error"))
	}
	p.web.RegisterServiceAPI(p)
}

func (p *power) deinitSvcWeb() {
	p.web.UnregisterServiceAPI(p)
}

func (p *power) WebRegisterRoutes(group *echo.Group) {
	group.GET("/poweroff", func(c echo.Context) error {
		data := &api.TokenRequest{}
		if !api.RequestLoader(c, data) {
			return nil
		}
		if err, uae := p.Poweroff(data.Token); err != nil {
			code := 500
			if uae {
				code = 403
			}
			api.ErrorResponseWriter(c, code, err)
		} else {
			api.ResponseWriter(c, nil)
		}
		return nil
	})
	group.GET("/reboot", func(c echo.Context) error {
		data := &api.TokenRequest{}
		if !api.RequestLoader(c, data) {
			return nil
		}
		if err, uae := p.Reboot(data.Token); err != nil {
			code := 500
			if uae {
				code = 403
			}
			api.ErrorResponseWriter(c, code, err)
		} else {
			api.ResponseWriter(c, nil)
		}
		return nil
	})
	group.POST("/wakealarm", func(c echo.Context) error {
		data := &wakeAlarmRequest{}
		if !api.RequestLoader(c, data) {
			return nil
		}
		if err, uae := p.SetWakeAlarm(time.Unix(data.Time, 0), data.Token); err != nil {
			code := 500
			if uae {
				code = 403
			}
			api.ErrorResponseWriter(c, code, err)
		} else {
			api.ResponseWriter(c, nil)
		}
		return nil
	})
	group.DELETE("/wakealarm", func(c echo.Context) error {
		data := &api.TokenRequest{}
		if !api.RequestLoader(c, data) {
			return nil
		}
		if err, uae := p.UnsetWakeAlarm(data.Token); err != nil {
			code := 500
			if uae {
				code = 403
			}
			api.ErrorResponseWriter(c, code, err)
		} else {
			api.ResponseWriter(c, nil)
		}
		return nil
	})
}
