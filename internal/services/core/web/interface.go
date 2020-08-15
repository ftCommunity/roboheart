package web

import "github.com/ftCommunity/roboheart/internal/service"

type Web interface {
	RegisterServiceAPI(webservice)
	UnregisterServiceAPI(service.Service)
}
