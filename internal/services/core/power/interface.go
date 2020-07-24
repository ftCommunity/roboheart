package power

import (
	"time"
)

type Power interface {
	Poweroff(token string) (error, bool)
	Reboot(token string) (error, bool)
	SetWakeAlarm(t time.Time, token string) (error, bool)
	UnsetWakeAlarm(token string) (error, bool)
}
