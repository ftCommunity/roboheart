package acm

import (
	"errors"
	"time"

	"github.com/thoas/go-funk"
)

type token struct {
	acm                           acm
	created, lastrefresh          time.Time
	totallifetime, refeshlifetime time.Duration
	permissions                   *map[string]bool
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
	if time.Now().Sub(t.created) > t.totallifetime {
		return false
	}
	if time.Now().Sub(t.lastrefresh) > t.refeshlifetime {
		return false
	}
	return true
}

func (t *token) Refresh() { t.lastrefresh = time.Now() }
