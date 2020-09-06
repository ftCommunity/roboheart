package testhelper

import (
	"github.com/ftCommunity-roboheart/roboheart/package/service"
	"testing"
)

func GetErrorFunc(t *testing.T) service.ErrorFunc {
	return func(e ...interface{}) {
		t.Error(e...)
	}
}
