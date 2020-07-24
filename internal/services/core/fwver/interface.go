package fwver

import (
	"github.com/blang/semver"
)

type FWVer interface {
	Get() semver.Version
	GetString() string
}
