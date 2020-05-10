package acm

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ftCommunity/roboheart/internal/service"
	"github.com/ftCommunity/roboheart/package/threadmanager"
	"github.com/thoas/go-funk"
)

var (
	DEFAULTS           []string = []string{"root", "user", "app"}
	TokenNotFoundError          = errors.New("Token not found")
)

const (
	tokenTotalLifetime   = 1 * time.Hour
	tokenRefreshLifetime = 15 * time.Minute
)

type acm struct {
	logger      service.LoggerFunc
	error       service.ErrorFunc
	permissions []string
	defaults    map[string]*map[string]bool
	tokens      map[string]*token
	tm          *threadmanager.ThreadManager
}

type ACM interface {
	RegisterPermission(name string, defaults map[string]bool) error
	CreateToken(defaults []string, layers ...map[string]bool) (string, error)
	UpdateToken(id string, layers ...map[string]bool) error
	GetToken(id string) (*token, error)
}

func (a *acm) Init(services map[string]service.Service, logger service.LoggerFunc, e service.ErrorFunc) error {
	a.logger = logger
	a.error = e
	a.permissions = make([]string, 0)
	a.defaults = make(map[string]*map[string]bool)
	for _, d := range DEFAULTS {
		a.addDefault(d)
	}
	a.tokens = make(map[string]*token)
	a.tm = threadmanager.NewThreadManager(a.logger, a.error)
	a.tm.Load("cleanup", a.cleanupThread)
	a.tm.Start("cleanup")
	return nil
}

func (a *acm) Stop() error                                                { a.tm.StopAll(); return nil }
func (a *acm) Name() string                                               { return "acm" }
func (a *acm) Dependencies() ([]string, []string)                         { return []string{}, []string{} }
func (a *acm) SetAdditionalDependencies(map[string]service.Service) error { return nil }
func (a *acm) UnsetAdditionalDependencies()                               {}

func (a *acm) cleanupThread(logger service.LoggerFunc, e service.ErrorFunc, stop, stopped chan interface{}) {
	for {
		select {
		case <-stop:
			{
				stopped <- struct{}{}
				return
			}
		case <-time.After(5 * time.Second):
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

func (a *acm) RegisterPermission(name string, defaults map[string]bool) error {
	if funk.ContainsString(a.permissions, name) {
		return errors.New("Permission already registered")
	}
	a.permissions = append(a.permissions, name)
	if defaults == nil {
		return nil
	}
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
	for _, l := range layers {
		for pn, ps := range l {
			if !funk.Contains(a.permissions, pn) {
				return errors.New("Permission not found")
			}
			if ps {
				if t.parent != nil {
					if !(*t.parent.permissions)[pn] {
						return errors.New("Cannot allow permission that is not allowed in parent")
					}
				}
			}
			(*t.permissions)[pn] = ps
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

func (a *acm) createToken() string {
	id, err := uuid.NewRandom()
	if err != nil {
		a.logger(err)
		return ""
	}
	t := new(token)
	t.totallifetime = tokenTotalLifetime
	t.refeshlifetime = tokenRefreshLifetime
	t.created = time.Now()
	t.Refresh()
	t.permissions = &map[string]bool{}
	a.tokens[id.String()] = t
	return id.String()
}

func (a *acm) addDefault(d string) {
	a.defaults[d] = &map[string]bool{}

}

var Service = new(acm)
