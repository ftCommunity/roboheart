package main

var services = map[string][3]string{
	//core
	"acm": {
		"github.com/ftCommunity/roboheart/internal/services/core/acm",
		"acm",
		"Service",
	},
	"config": {
		"github.com/ftCommunity/roboheart/internal/services/core/config",
		"config",
		"Service",
	},
	"deviceinfo": {
		"github.com/ftCommunity/roboheart/internal/services/core/deviceinfo",
		"deviceinfo",
		"Service",
	},
	"filesystem": {
		"github.com/ftCommunity/roboheart/internal/services/core/filesystem",
		"filesystem",
		"Service",
	},
	"fwver": {
		"github.com/ftCommunity/roboheart/internal/services/core/fwver",
		"fwver",
		"Service",
	},
	"locale": {
		"github.com/ftCommunity/roboheart/internal/services/core/locale",
		"locale",
		"Service",
	},
	"power": {
		"github.com/ftCommunity/roboheart/internal/services/core/power",
		"power",
		"Service",
	},
	"web": {
		"github.com/ftCommunity/roboheart/internal/services/core/web",
		"web",
		"Service",
	},

	//additional
	"pkgmanager": {
		"github.com/ftCommunity/roboheart/internal/services/pkgmanager",
		"pkgmanager",
		"Service",
	},
	"releasever": {
		"github.com/ftCommunity/roboheart/internal/services/releasever",
		"relver",
		"Service",
	},
	"vncserver": {
		"github.com/ftCommunity/roboheart/internal/services/vncserver",
		"vncserver",
		"Service",
	},
}
