package web

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/gorilla/mux"
)

type Web interface {
	RegisterServiceAPI(service.Service) *mux.Router
}
