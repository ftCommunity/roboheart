package acm

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/testhelper"
	"testing"
)

type testenv struct {
	acm ACM
}

func initTestEnv(t *testing.T) *testenv {
	env := &testenv{}
	acm := &acm{}
	if err := acm.Init(map[string]service.Service{}, func(e ...interface{}) {}, testhelper.GetErrorFunc(t)); err != nil {
		t.Error(err)
	}
	env.acm = acm
	return env
}
