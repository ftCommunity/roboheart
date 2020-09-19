module github.com/ftCommunity-roboheart/roboheart

go 1.14

replace github.com/digineo/go-uci => github.com/ftCommunity-roboheart/go-uci v0.0.0-20200725220005-3826098d8ac7

require (
	github.com/akamensky/argparse v1.2.2
	github.com/blang/semver v3.5.1+incompatible
	github.com/ftCommunity-roboheart/roboheart-svc-releasever v0.0.0-20200913123916-c95d25e8d587
	github.com/ftCommunity-roboheart/roboheart-svc-vncserver v0.0.0-20200913123943-a19f7c20f49a
	github.com/ftCommunity-roboheart/roboheart-svcs-core v0.0.0-20200919215026-5e786780b496
	github.com/ftCommunity-roboheart/roboheart-svcs-net v0.0.0-20200919215513-e2780e12119c
	github.com/google/uuid v1.1.2
	github.com/labstack/echo/v4 v4.1.17
	github.com/spf13/afero v1.4.0
	github.com/thoas/go-funk v0.7.0
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	golang.org/x/sys v0.0.0-20200909081042-eff7692f9009 // indirect
)
