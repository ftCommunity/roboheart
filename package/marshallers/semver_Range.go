package marshallers

import (
	"github.com/blang/semver"
)

type Range struct {
	semver.Range
}

func (r *Range) MarshalJSON() ([]byte, error) {
	return MakeByteString(""), nil
}

func (r *Range) UnmarshalJSON(data []byte) error {
	var err error
	r.Range, err = semver.ParseRange(extractString(data))
	return err
}
