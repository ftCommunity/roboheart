package acm

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/internal/services/core/robotime"
	"github.com/ftCommunity/roboheart/package/testhelper"
	"testing"
)

type testenv struct {
	acm ACM
	rt *robotime.RoboTimeMock
}

func initTestEnv(t *testing.T) *testenv {
	env := &testenv{}
	acm := &acm{}
	var rt service.Service
	rt= robotime.NewMock()
	env.rt=rt.(*robotime.RoboTimeMock)
	acm.Init(map[string]service.Service{"robotime": rt}, func(e ...interface{}) {}, testhelper.GetErrorFunc(t))
	env.acm = acm
	return env
}
