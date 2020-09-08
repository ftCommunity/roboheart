module github.com/ftCommunity-roboheart/roboheart

go 1.14

replace github.com/digineo/go-uci => github.com/ftCommunity-roboheart/go-uci v0.0.0-20200725220005-3826098d8ac7

require (
	github.com/akamensky/argparse v1.2.2
	github.com/blang/semver v3.5.1+incompatible
	github.com/ftCommunity-roboheart/roboheart-svc-releasever v0.0.0-20200906133623-79b2f6c79149
	github.com/ftCommunity-roboheart/roboheart-svc-vncserver v0.0.0-20200906135251-445e1885fd1e
	github.com/ftCommunity-roboheart/roboheart-svcs-core v0.0.0-20200906114943-65a04c29b345
	github.com/google/uuid v1.1.2
	github.com/labstack/echo/v4 v4.1.17
	github.com/spf13/afero v1.3.5
	github.com/thoas/go-funk v0.7.0
)
