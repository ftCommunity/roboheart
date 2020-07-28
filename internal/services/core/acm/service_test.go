package acm

import (
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"testing"
)

func TestAcm_Init(t *testing.T) {
	initTestEnv(t)
}

func TestAcm_RegisterPermissions_OK(t *testing.T) {
	env := initTestEnv(t)
	for i := 0; i < 10; i++ {
		rid, err := uuid.NewRandom()
		if err != nil {
			t.Error(err)
		}
		pid := rid.String()
		defaults := map[string]bool{}
		for di, d := range DEFAULTS {
			defaults[d] = (i+di)%2 == 0
		}
		if err := env.acm.RegisterPermission(pid, defaults, map[string]string{}); err != nil {
			t.Error(err)
		}
	}
}

func TestAcm_RegisterPermissions_DoubleRegistration(t *testing.T) {
	env := initTestEnv(t)
	rid, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	pid := rid.String()
	if err := env.acm.RegisterPermission(pid, map[string]bool{}, map[string]string{}); err != nil {
		t.Error(err)
	}
	if err := env.acm.RegisterPermission(pid, map[string]bool{}, map[string]string{}); err == nil {
		t.Error("Permission was registered twice")
	}
}

func TestAcm_RegisterPermission_UnknownDefault(t *testing.T) {
	env := initTestEnv(t)
	rid, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	pid := rid.String()
	var d string
defaultfinder:
	for {
		rd, err := uuid.NewRandom()
		if err != nil {
			t.Error(err)
		}
		d = rd.String()
		if !funk.ContainsString(DEFAULTS, d) {
			break defaultfinder
		}
	}
	if err := env.acm.RegisterPermission(pid, map[string]bool{d: true}, map[string]string{}); err == nil {
		t.Error("Used unknown default")
	}
}
