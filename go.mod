module github.com/ftCommunity-roboheart/roboheart

go 1.14

replace github.com/digineo/go-uci => github.com/ftCommunity-roboheart/go-uci v0.0.0-20200725220005-3826098d8ac7

require (
	github.com/akamensky/argparse v1.2.2
	github.com/blang/semver v3.5.1+incompatible
	github.com/ftCommunity-roboheart/roboheart-svcs-core v0.0.0-20201114140116-16b44e0c6d7b
	github.com/google/uuid v1.1.2
	github.com/labstack/echo/v4 v4.1.17
	github.com/spf13/afero v1.4.1
	github.com/stretchr/testify v1.6.1
	github.com/thoas/go-funk v0.7.0
	golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c // indirect
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb // indirect
	golang.org/x/sys v0.0.0-20201204225414-ed752295db88 // indirect
)
