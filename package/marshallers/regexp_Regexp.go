package marshallers

import (
	"regexp"
	"strings"
)

type Regexp struct {
	regexp.Regexp
}

func (r Regexp) MarshalJSON() ([]byte, error) {
	return MakeByteString(strings.TrimSuffix(strings.TrimPrefix(r.String(), "^"), "$")), nil
}

func (r *Regexp) UnmarshalJSON(data []byte) error {
	var err error
	reg, err := regexp.Compile("^" + extractString(data) + "$")
	r.Regexp = *reg
	return err
}
