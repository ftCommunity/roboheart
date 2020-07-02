package marshallers

import (
	"regexp"
)

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) MarshalJSON() ([]byte, error) {
	return MakeByteString(r.String()), nil
}

func (r *Regexp) UnmarshalJSON(data []byte) error {
	var err error
	r.Regexp, err = regexp.Compile(extractString(data))
	return err
}
