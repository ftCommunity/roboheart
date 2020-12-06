package marshallers

import (
	"bytes"
	"github.com/google/uuid"
	"testing"
)

func TestUUID_MarshalJSON(t *testing.T) {
	for i := 0; i < 10; i++ {
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatal(err)
		}
		us := UUID(id)
		j, err := us.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(j, []byte("\""+id.String()+"\"")) {
			t.Error("Output mismatch")
		}
	}
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	for i := 0; i < 10; i++ {
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatal(err)
		}
		raw := "\"" + id.String() + "\""
		us := UUID{}
		err = us.UnmarshalJSON([]byte(raw))
		if err != nil {
			t.Fatal(err)
		}
		if id != uuid.UUID(us) {
			t.Fatal("UUID mismatch")
		}
	}
	us := UUID{}
	if us.UnmarshalJSON([]byte("abc")) == nil {
		t.Fatal("Length mismatch not detected")
	}

}
