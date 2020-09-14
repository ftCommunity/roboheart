package services

import (
	"github.com/ftCommunity-roboheart/roboheart-svc-releasever/package/services/releasever"
	"github.com/ftCommunity-roboheart/roboheart-svc-vncserver/package/services/vncserver"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/acm"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/config"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/deviceinfo"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/filesystem"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/fwver"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/locale"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/power"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/robotime"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/serviceadmin"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services/core/web"
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services/doh"
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services/net"
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services/net-svcs/discovery"
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services/net-svcs/mesh/alfred"
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services/net-svcs/mesh/batman"
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services/net-svcs/mesh/vis"
	"github.com/ftCommunity-roboheart/roboheart/package/service"
)

var Services = []service.Service{
	acm.Service,
	alfred.Service,
	batman.Service,
	config.Service,
	deviceinfo.Service,
	discovery.Service,
	doh.Service,
	filesystem.Service,
	fwver.Service,
	locale.Service,
	net.Service,
	power.Service,
	relver.Service,
	robotime.Service,
	serviceadmin.Service,
	vis.Service,
	vncserver.Service,
	web.Service,
}
