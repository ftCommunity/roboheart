package serviceadmin

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/labstack/echo/v4"
)

func (s *serviceadmin) initSvcWeb(svc service.Service) {
	var ok bool
	s.web, ok = svc.(web.Web)
	if !ok {
		s.error(errors.New("Type assertion error"))
	}
	s.web.RegisterServiceAPI(s)
}

func (s *serviceadmin) deinitSvcWeb() {
	s.web.UnregisterServiceAPI(s)
}

func (s *serviceadmin) WebRegisterRoutes(group *echo.Group) {
	group.GET("/services", func(c echo.Context) error {
		api.CheckACMAPICall(c, s.acm, READ_PERM, func() {
			api.ResponseWriter(c, s.manager.GetServicesInfo())
		})
		return nil
	})
	group.GET("/service/:sn", func(c echo.Context) error {
		api.CheckACMAPICall(c, s.acm, READ_PERM, func() {
			if info, err := s.manager.GetServiceInfo(c.Param("sn")); err != nil {
				api.ErrorResponseWriter(c, 404, err)
			} else {
				api.ResponseWriter(c, info)
			}
		})
		return nil
	})
}
