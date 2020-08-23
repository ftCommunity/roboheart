package acm

import (
	"errors"
	"github.com/ftCommunity/roboheart/internal/services/core/robotime"
	"github.com/ftCommunity/roboheart/package/servicehelpers"
	"time"

	"github.com/google/uuid"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/threadmanager"
)

type acm struct {
	logger      service.LoggerFunc
	error       service.ErrorFunc
	rt robotime.RoboTime
	permissions map[string]map[string]string
	defaults    map[string]*map[string]bool
	tokens      map[string]*token
	tm          *threadmanager.ThreadManager
}

func (a *acm) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) {
	a.logger = logger
	a.error = e
	if err := servicehelpers.CheckMainDependencies(a, services); err != nil {
		e(err)
	}
	if err := servicehelpers.InitializeDependencies(services, servicehelpers.ServiceInitializers{"robotime": a.initSvcRoboTime}); err != nil {
		e(err)
	}
	a.permissions = make(map[string]map[string]string)
	a.defaults = make(map[string]*map[string]bool)
	for _, d := range DEFAULTS {
		a.addDefault(d)
	}
	a.tokens = make(map[string]*token)
	a.tm = threadmanager.NewThreadManager(a.logger, a.error)
	a.tm.Load("cleanup", a.cleanupThread)
	a.tm.Start("cleanup")
}

func (a *acm) Stop()        { a.tm.StopAll() }
func (a *acm) Name() string { return "acm" }
func (a *acm) Dependencies() service.ServiceDependencies {
	return service.ServiceDependencies{
		Deps: []string{"robotime"},
	}
}

func (a *acm) cleanupThread(_ service.LoggerFunc, _ service.ErrorFunc, stop, stopped chan interface{}) {
	for {
		select {
		case <-stop:
			{
				stopped <- struct{}{}
				return
			}
		case <-a.rt.After(5 * time.Second):
			{
				for id, t := range a.tokens {
					if !t.CheckValid() {
						hasChild := false
						for _, tt := range a.tokens {
							if tt.parent == t {
								hasChild = true
							}
						}
						if !hasChild {
							delete(a.tokens, id)
						}
					}
				}
			}
		}
	}
}

func (a *acm) RegisterPermission(name string, defaults map[string]bool, desc map[string]string) error {
	if _, ok := a.permissions[name]; ok {
		return errors.New("Permission already registered")
	}
	a.permissions[name] = desc
	if defaults == nil {
		return nil
	}
	(*a.defaults["root"])[name] = true
	for dn, ds := range defaults {
		if _, ok := a.defaults[dn]; !ok {
			return errors.New("Unknown default " + dn)
		}
		(*a.defaults[dn])[name] = ds
	}
	return nil
}

func (a *acm) getDafault(name string) (map[string]bool, error) {
	if d, ok := a.defaults[name]; ok {
		return *d, nil
	}
	return nil, errors.New("Default not found")
}

func (a *acm) CreateToken(defaults []string, layers ...map[string]bool) (string, error) {
	id := a.createToken()
	defaultlayers := []map[string]bool{}
	for _, dn := range defaults {
		d, err := a.getDafault(dn)
		if err != nil {
			return "", err
		}
		defaultlayers = append(defaultlayers, d)
	}
	if err := a.UpdateToken(id, append(defaultlayers, layers...)...); err != nil {
		return "", err
	}
	return id, nil
}

func (a *acm) UpdateToken(id string, layers ...map[string]bool) error {
	t, ok := a.tokens[id]
	if !ok {
		return errors.New("Token not found")
	}
	disabled := map[string]bool{}
	for _, l := range layers {
		for pn, ps := range l {
			if _, ok := a.permissions[pn]; !ok {
				return errors.New("Permission not found")
			}
			if ps {
				if t.parent != nil {
					if !(*t.parent.permissions)[pn] {
						return errors.New("Cannot allow permission that is not allowed in parent")
					}
				}
			} else {
				disabled[pn] = false
			}
			(*t.permissions)[pn] = ps
		}
	}
	for _, t := range a.tokens {
		if t.parent == t {
			a.UpdateToken(t.id, disabled)
		}
	}
	return nil
}

func (a *acm) GetToken(id string) (*token, error) {
	t, ok := a.tokens[id]
	if !ok {
		return nil, TokenNotFoundError
	}
	return t, nil
}

func (a *acm) CheckTokenPermission(token string, permission string) (error, bool) {
	t, err := a.GetToken(token)
	if err != nil {
		return err, isUserError(err)
	}
	t.Refresh()
	perm, err := t.GetPermission(permission)
	if err != nil {
		return err, isUserError(err)
	}
	if !perm {
		return NotPermittedError, true
	}
	return nil, false
}

func (a *acm) GetPermissionDescription(name string) (map[string]string, error) {
	if desc, ok := a.permissions[name]; ok {
		return desc, nil
	} else {
		return nil, errors.New("Permission not found")
	}
}

func (a *acm) createToken() string {
	id, err := uuid.NewRandom()
	if err != nil {
		a.logger(err)
		return ""
	}
	t := new(token)
	t.refeshlifetime = tokenRefreshLifetime
	t.Refresh()
	t.permissions = &map[string]bool{}
	t.id = id.String()
	a.tokens[id.String()] = t
	return id.String()
}

func (a *acm) addDefault(d string) {
	a.defaults[d] = &map[string]bool{}

}
