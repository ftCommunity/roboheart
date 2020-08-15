package web

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/labstack/echo/v4"
)

type webservice interface {
	service.Service
	WebRegisterRoutes(group *echo.Group)
}
