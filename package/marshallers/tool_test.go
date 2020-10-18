package marshallers

import (
	"bytes"
	"testing"
)

func TestExtractString(t *testing.T) {
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
		/*0*/ "abc",
		/*1*/ "abc",
		/*2*/ "abc",
		/*3*/ "abc",
		/*4*/ "",
		/*5*/ "",
		/*6*/ " ",
		/*7*/ " ",
		/*8*/ "",
	}
	for i, ind := range in {
		if result := extractString(ind); result != out[i] {
			t.Error("Output mismatch")
		}
	}
}

func TestMakeByteString(t *testing.T) {
	if !bytes.Equal(MakeByteString("abc"), []byte{'"', 'a', 'b', 'c', '"'}) {
		t.Error("Output mismatch")
	}
}
