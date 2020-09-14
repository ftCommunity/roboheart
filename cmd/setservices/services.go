package main

import (
	"github.com/ftCommunity-roboheart/roboheart-svc-releasever/package/services"
	"github.com/ftCommunity-roboheart/roboheart-svc-vncserver/package/services"
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services"
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services"
)

var serviceproviders = map[string]map[string][3]string{
	"github.com/ftCommunity-roboheart/roboheart-svcs-core/package/services":      svcs_core.Services,
	"github.com/ftCommunity-roboheart/roboheart-svc-releasever/package/services": svc_releasever.Services,
	"github.com/ftCommunity-roboheart/roboheart-svc-vncserver/package/services":  svc_vncserver.Services,
	"github.com/ftCommunity-roboheart/roboheart-svcs-net/package/services":       svcs_net.Services,
}
