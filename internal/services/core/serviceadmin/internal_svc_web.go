package serviceadmin

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
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

}
