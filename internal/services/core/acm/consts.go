package acm

import (
	"time"
)

const (
	tokenRefreshLifetime = 15 * time.Minute
)

var (
	DEFAULTS = []string{"root", "user", "app"}
)
