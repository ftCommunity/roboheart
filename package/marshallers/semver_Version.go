package marshallers

import (
	"github.com/blang/semver"
)

type Version struct {
	semver.Version
}

func (v Version) MarshalJSON() ([]byte, error) {
	return MakeByteString(v.String()), nil
}

func (v *Version) UnmarshalJSON(data []byte) error {
	var err error
	ver, err := semver.New(extractString(data))
	v.Version = *ver
	return err
}
