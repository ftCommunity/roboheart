# roboheart
roboheart is the new central daemon for the [ftCommunity Firmware](https://github.com/ftCommunity/ftcommunity-TXT). Its main goal is to provide a central place where all configuration and system access takes place while having a maximum of security by using a strong ACL system.

## Usage
### Start-up
`./roboheart` ;-)
### Stop
`ctrl-c`

## Building
### Download sources
```
go get github.com/ftCommunity/roboheart
cd $HOME/go/src/github.com/ftCommunity/roboheart
go mod vendor
```

### Build for TXT
`GOARCH=arm go build -o roboheart cmd/roboheart/main.go`

## API
You can test the web API by running `curl <TXT-IP>:8080/api/fwver/version`
