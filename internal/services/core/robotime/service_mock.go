package robotime

import (
	"sync"
	"time"
)

type RoboTimeMock struct {
	cTime     time.Time
	timers    map[time.Duration][]chan time.Time
	timerlock sync.Mutex
}


func (r RoboTimeMock) Now() time.Time {
	return r.cTime
}

func (r *RoboTimeMock) SetTime(t time.Time) {
	r.cTime = t
}

func (r *RoboTimeMock) After(d time.Duration) <-chan time.Time {
	r.timerlock.Lock()
	c := make(chan time.Time)
	if _, ok := r.timers[d]; !ok {
		r.timers[d] = []chan time.Time{}
	}
	r.timers[d] = append(r.timers[d], c)
	r.timerlock.Unlock()
	return c
}

func (r *RoboTimeMock) Fire(d time.Duration) {
	r.timerlock.Lock()
	defer r.timerlock.Unlock()
	if tl, ok := r.timers[d]; ok {
		for _, t := range tl {
			t <- r.cTime
		}
		delete(r.timers, d)
	}
}

func NewMock() *RoboTimeMock {
	return &RoboTimeMock{
		timers: map[time.Duration][]chan time.Time{},
	}
}
