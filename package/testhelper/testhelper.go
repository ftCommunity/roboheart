package testhelper

import (
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
	"testing"
)

func GetErrorFunc(t *testing.T) instance.ErrorFunc {
	return func(e ...interface{}) {
		t.Error(e...)
	}
}
