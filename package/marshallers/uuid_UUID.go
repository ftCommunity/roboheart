package marshallers

import (
	"github.com/google/uuid"
)

type UUID [16]byte

func (u UUID) MarshalJSON() ([]byte, error) {
	return MakeByteString(uuid.UUID(u).String()), nil
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	id, err := uuid.Parse(extractString(data))
	*u=UUID(id)
	return err
}

func (u UUID) MarshalText() ([]byte, error) {
	return u.MarshalJSON()
}

func (u *UUID) UnmarshalText(text []byte) error {
	return u.UnmarshalJSON(text)
}