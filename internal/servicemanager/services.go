package servicemanager

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/ftCommunity/roboheart/internal/services/core/config"
	"github.com/ftCommunity/roboheart/internal/services/core/deviceinfo"
	"github.com/ftCommunity/roboheart/internal/services/core/filesystem"
	"github.com/ftCommunity/roboheart/internal/services/core/fwver"
	"github.com/ftCommunity/roboheart/internal/services/core/locale"
	"github.com/ftCommunity/roboheart/internal/services/core/power"
	"github.com/ftCommunity/roboheart/internal/services/core/web"
	"github.com/ftCommunity/roboheart/internal/services/pkgmanager"
	"github.com/ftCommunity/roboheart/internal/services/releasever"
	"github.com/ftCommunity/roboheart/internal/services/vncserver"
)

var services = []service.Service{
	acm.Service,
	config.Service,
	deviceinfo.Service,
	filesystem.Service,
	fwver.Service,
	locale.Service,
	pkgmanager.Service,
	power.Service,
	relver.Service,
	vncserver.Service,
	web.Service,
}
