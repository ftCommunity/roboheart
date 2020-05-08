package fwver

import (
	"io/ioutil"
	"strings"

	"github.com/blang/semver"
	"github.com/ftCommunity/roboheart/internal/service"
)

type fwver struct {
	rawver string
	semver semver.Version
}

type FWVer interface {
	Get() semver.Version
	GetString() string
}

func (f *fwver) Init(map[string]service.Service, service.LoggerFunc, service.ErrorFunc) error {
	raw, err := ioutil.ReadFile("/etc/fw-ver.txt")
	if err != nil {
		return err
	}
	f.rawver = strings.Split(string(raw), "\n")[0]
	f.semver, err = semver.Make(f.rawver)
	if err != nil {
		return err
	}
	return nil
}

func (f *fwver) Stop() error                                                { return nil }
func (f *fwver) Name() string                                               { return "fwver" }
func (f *fwver) Dependencies() ([]string, []string)                         { return []string{}, []string{} }
func (f *fwver) SetAdditionalDependencies(map[string]service.Service) error { return nil }
func (f *fwver) UnsetAdditionalDependencies()                               {}

func (f *fwver) Get() semver.Version {
	return f.semver
}

func (f *fwver) GetString() string {
	return f.rawver
}

var Service = new(fwver)
