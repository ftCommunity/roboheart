package web

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/api"
	"github.com/gorilla/mux"
)

type web struct {
	logger      service.LoggerFunc
	error       service.ErrorFunc
	mux, apiMux *mux.Router
	services    []string
	srv         *http.Server
	srvwg       sync.WaitGroup
}

type Web interface {
	RegisterServiceAPI(service.Service) *mux.Router
}

func (w *web) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	w.logger = logger
	w.error = e
	w.mux = mux.NewRouter()
	w.mux.Use(serverHeaderMiddleware)
	w.mux.NotFoundHandler = w.mux.NewRoute().BuildOnly().HandlerFunc(http.NotFound).GetHandler()
	w.mux.MethodNotAllowedHandler = w.mux.NewRoute().BuildOnly().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}).GetHandler()
	w.apiMux = getSubMux(w.mux, "/api")
	w.apiMux.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	w.apiMux.Use(corsHeadersMiddleware)
	w.apiMux.NotFoundHandler = w.apiMux.NewRoute().BuildOnly().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ErrorResponseWriter(w, 404, errors.New("Page not found"))
	}).GetHandler()
	w.apiMux.MethodNotAllowedHandler = w.apiMux.NewRoute().BuildOnly().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ErrorResponseWriter(w, 405, errors.New("Method not allowed"))
	}).GetHandler()
	w.srv = &http.Server{Addr: ":8080", Handler: w.mux}
	go func() {
		w.srvwg.Add(1)
		defer w.srvwg.Done()
		if err := w.srv.ListenAndServe(); err != http.ErrServerClosed {
			w.error(err)
		}
	}()
	return nil
}

func (w *web) Stop() error {
	if err := w.srv.Shutdown(context.TODO()); err != nil {
		return err
	}
	w.srvwg.Wait()
	return nil
}

func (w *web) Name() string                                               { return "web" }
func (w *web) Dependencies() ([]string, []string)                         { return []string{}, []string{} }
func (w *web) SetAdditionalDependencies(map[string]service.Service) error { return nil }
func (w *web) UnsetAdditionalDependencies()                               {}

func (w *web) RegisterServiceAPI(s service.Service) *mux.Router {
	m := getSubMux(w.apiMux, "/"+s.Name())
	w.services = append(w.services, s.Name())
	return m
}

func getSubMux(m *mux.Router, p string) *mux.Router {
	return m.PathPrefix(p).Subrouter()
}

var Service = new(web)
