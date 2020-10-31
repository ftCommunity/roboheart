package services

import (
	svcs_core "github.com/ftCommunity-roboheart/roboheart-svcs-core/package"
	"github.com/ftCommunity-roboheart/roboheart/package/manifest"
)

var ServiceProviders = [][]manifest.ServiceManifest{
	svcs_core.Services,
}
