package services

import (
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
	"github.com/ftCommunity-roboheart/roboheart/package/service"
)

var Services = []service.Service{
	acm.Service,
	config.Service,
	deviceinfo.Service,
	filesystem.Service,
	fwver.Service,
	locale.Service,
	power.Service,
	robotime.Service,
	serviceadmin.Service,
	web.Service,
}
