package marshallers

import (
	"bytes"
	"regexp"
	"testing"
)

func TestRegexp_MarshalJSON(t *testing.T) {
	for _, raw := range []string{"abc", "^abc", "abc$", "^abc$"} {
		r, err := regexp.Compile(raw)
		if err != nil {
			t.Error(err)
		}
		rs := &Regexp{Regexp: *r}
		j, err := rs.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(j, []byte{'"', 'a', 'b', 'c', '"'}) {
			t.Error("Output mismatch")
		}
	}
}

func TestRegexp_UnmarshalJSON(t *testing.T) {
	in := [][]byte{
		/*0*/ {'a', 'b', 'c'},
		/*1*/ {'"', 'a', 'b', 'c', '"'},
		/*2*/ {'"', 'a', 'b', 'c'},
		/*3*/ {'a', 'b', 'c', '"'},
		/*4*/ {'"', '"'},
		/*5*/ {'"'},
		/*6*/ {' ', '"'},
		/*7*/ {'"', ' '},
		/*8*/ {},
	}
	out := []string{
		/*0*/ "^abc$",
		/*1*/ "^abc$",
		/*2*/ "^abc$",
		/*3*/ "^abc$",
		/*4*/ "^$",
		/*5*/ "^$",
		/*6*/ "^ $",
		/*7*/ "^ $",
		/*8*/ "^$",
	}
	for i, ind := range in {
		rs := &Regexp{}
		if err := rs.UnmarshalJSON(ind); err != nil {
			t.Error(err)
		}
		if rs.Regexp.String() != out[i] {
			t.Error("Output mismatch")
		}
	}
}
