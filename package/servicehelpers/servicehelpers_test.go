package servicehelpers

import (
	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/google/uuid"
	"testing"
)

func TestCheckDependencies(t *testing.T) {
	for i := 0; i < 10; i++ {
		var req []string
		services := map[string]service.Service{}
		for j := 0; j < 5; j++ {
			suuid, err := uuid.NewRandom()
			if err != nil {
				t.Error(err)
			}
			sid := suuid.String()
			req = append(req, sid)
			services[sid] = nil
		}
		if err := CheckMainDependencies(&dummyService{}, services); err != nil {
			t.Error(err)
		}
		if err := CheckMainDependencies(&dummyService{}, services); err != nil {
			t.Error(err)
		}
		if err := CheckMainDependencies(&dummyService{deps: req}, map[string]service.Service{}); err == nil {
			t.Error("Data mismatch not detected")
		}
		if err := CheckAdditionalDependencies(&dummyService{adeps: req}, map[string]service.Service{}); err == nil {
			t.Error("Data mismatch not detected")
		}
		if err := CheckMainDependencies(&dummyService{deps: req}, services); err != nil {
			t.Error(err)
		}
		if err := CheckAdditionalDependencies(&dummyService{adeps: req}, services); err != nil {
			t.Error(err)
		}
	}

}
