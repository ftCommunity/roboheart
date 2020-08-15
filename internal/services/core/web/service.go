package web

import (
	"context"
	"errors"
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
)

type web struct {
	logger   service.LoggerFunc
	error    service.ErrorFunc
	echo     *echo.Echo
	services map[string]webservice
	lock     sync.Mutex
}

func (w *web) Init(_ map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) {
	w.logger = logger
	w.error = e
	w.services = make(map[string]webservice)
	w.echo = echo.New()
	w.echo.Logger.SetLevel(log.OFF)
	w.echo.HideBanner = true
	w.echo.HidePort = true
	w.start()
}

func (w *web) Stop() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.stop()
}

func (w *web) Name() string { return "web" }

func (w *web) RegisterServiceAPI(svc webservice) {
	if _, ok := w.services[svc.Name()]; ok {
		w.error(errors.New(svc.Name() + " was registered twice"))
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	w.services[svc.Name()] = svc
	svc.WebRegisterRoutes(w.getServiceGroup(svc.Name()))
}

func (w *web) UnregisterServiceAPI(svc service.Service) {
	if _, ok := w.services[svc.Name()]; !ok {
		w.error(errors.New(svc.Name() + " is not registered"))
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	delete(w.services, svc.Name())
	w.restart()
}

func (w *web) start() {
	go func() {
		if err := w.echo.Start(":8080"); err != nil {
			if err != http.ErrServerClosed {
				w.error(err)
			}
		}
	}()
}

func (w *web) stop() {
	if err := w.echo.Shutdown(context.TODO()); err != nil {
		w.error(err)
	}
}

func (w *web) restart() {
	w.stop()
	w.start()
	for sn, svc := range w.services {
		svc.WebRegisterRoutes(w.getServiceGroup(sn))
	}
}

func (w *web) getServiceGroup(name string) *echo.Group {
	return w.echo.Group("/api/"+name, middleware.CORS())
}
