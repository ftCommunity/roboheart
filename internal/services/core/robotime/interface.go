package robotime

import "time"

type RoboTime interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}
