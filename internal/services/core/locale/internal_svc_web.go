package locale

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/labstack/echo/v4"
)

func (l *locale) initSvcWeb(svc service.Service) {
	var ok bool
	l.web, ok = svc.(web.Web)
	if !ok {
		l.error(errors.New("Type assertion error"))
	}
	l.web.RegisterServiceAPI(l)
}

func (l *locale) deinitSvcWeb() {
	l.web.UnregisterServiceAPI(l)
}

func (l *locale) WebRegisterRoutes(group *echo.Group) {
	group.GET("/locale", func(c echo.Context) error {
		if locale, err := l.GetLocale(); err != nil {
			api.ErrorResponseWriter(c, 500, err)
		} else {
			api.ResponseWriter(c, locale)
		}
		return nil
	})
	group.POST("/locale", func(c echo.Context) error {
		data := &struct {
			api.TokenRequest
			Locale string `json:"locale"`
		}{}
		if !api.RequestLoader(c, data) {
			return nil
		}
		if err, uae := l.SetLocale(data.Token, data.Locale); err != nil {
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
	group.GET("/allowed", func(c echo.Context) error {
		api.ResponseWriter(c, l.GetAllowedLocales())
		return nil
	})
}
