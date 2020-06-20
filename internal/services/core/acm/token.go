package acm

import (
	"errors"
	"time"

	"github.com/thoas/go-funk"
)

type token struct {
	acm            acm
	lastrefresh    time.Time
	refeshlifetime time.Duration
	permissions    *map[string]bool
	parent         *token
	id             string
}

func (t *token) GetPermission(name string) (bool, error) {
	if ps, ok := (*t.permissions)[name]; ok {
		return ps, nil
	}
	if funk.Contains(t.acm.permissions, name) {
		return false, nil
	}
	return false, errors.New("Permission not found")
}

func (t *token) CheckValid() bool {
	if time.Now().Sub(t.lastrefresh) > t.refeshlifetime {
		return false
	}
	return true
}

func (t *token) MakeSubToken(defaults []string, layers ...map[string]bool) (string, error) {
	id := t.makeSub()
	defaultlayers := []map[string]bool{}
	for _, dn := range defaults {
		d, err := t.acm.getDafault(dn)
		if err != nil {
			return "", err
		}
		defaultlayers = append(defaultlayers, d)
	}
	if err := t.acm.UpdateToken(id, append(defaultlayers, layers...)...); err != nil {
		return "", err
	}
	return id, nil
}

func (t *token) makeSub() string {
	s := t.acm.createToken()
	st, err := t.acm.GetToken(s)
	if err != nil {
		t.acm.error(err)
	}
	st.parent = t
	return s
}

func (t *token) Refresh() { t.lastrefresh = time.Now() }
