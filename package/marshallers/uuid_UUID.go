package marshallers

import (
	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func (u*UUID) MarshalJSON() ([]byte, error) {

	return MakeByteString(u.String()), nil
}

func (u*UUID) UnmarshalJSON(data []byte) error {
	var err error
	u.UUID, err = uuid.Parse(extractString(data))
	return err
}
