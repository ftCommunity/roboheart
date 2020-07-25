package relver

import "github.com/blang/semver"

type Release struct {
	Version  semver.Version `json:"version"`
	Download string         `json:"filename"`
	Size     int            `json:"size"`
}
