package servicemanager

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/ftCommunity/roboheart/internal/services/core/fwver"
	"github.com/ftCommunity/roboheart/internal/services/core/locale"
	"github.com/ftCommunity/roboheart/internal/services/core/power"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/internal/services/pkgmanager"
	"github.com/ftCommunity/roboheart/internal/services/releasever"
	"github.com/ftCommunity/roboheart/internal/services/vncserver"
)

var (
	services = map[string]service.Service{
		"acm":    acm.Service,
		"config": config.Service,
		"fwver":  fwver.Service,
		"locale": locale.Service,
		"relver": relver.Service,
		"power":  power.Service,
		"web":    web.Service,

		"pkgmanager": pkgmanager.Service,
		"vncserver":  vncserver.Service,
	}

	buildservices = map[string]bool{
		"acm":    true,
		"config": true,
		"fwver":  true,
		"locale": true,
		"relver": true,
		"power":  true,
		"web":    true,

		"pkgmanager": true,
		"vncserver":  true,
	}
)

func getServices() []service.Service {
	sl := []service.Service{}
	for sn, ss := range buildservices {
		if ss {
			sl = append(sl, services[sn])
		}
	}
	return sl
}
