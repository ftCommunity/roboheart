package robotime

import (
	"sync"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	r := NewMock()
	var i, j int64
	for i = 0; i < 1000; i += 100 {
		for j = 0; j < 1000; j += 100 {
			r.SetTime(time.Unix(i, j))
			if !r.Now().Equal(time.Unix(i, j)) {
				t.Error("Mismatch")
			}
		}
	}
}

func aftertest(t *testing.T, r RoboTime, d time.Duration, cond *sync.Cond, c *int, mode int) {
	running := make(chan interface{})
	go func() {
		done := make(chan struct{})
		go func() {
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
			close(done)
		}()
		running<- struct {}{}
		select {
		case <-r.After(d + time.Duration(mode)):
			{
				if mode == 0 {
					cond.L.Lock()
					*c--
					cond.L.Unlock()
				} else {
					t.Error("Mismatch")
				}
				return
			}
		case <-done:
			if mode == 0 {
				t.Error("Mismatch")
			} else {
				cond.L.Lock()
				*c--
				cond.L.Unlock()
			}
			return
		}
	}()
	<-running
}

func TestAfter(t *testing.T) {
	var d time.Duration
	for d = 0; d < 1000; d += 100 {
		r := NewMock()
		c:=9
		cond := sync.NewCond(&sync.Mutex{})
		aftertest(t, r, d, cond, &c, 0)
		aftertest(t, r, d, cond, &c, 0)
		aftertest(t, r, d, cond, &c, 0)
		aftertest(t, r, d, cond, &c, 1)
		aftertest(t, r, d, cond, &c, 2)
		aftertest(t, r, d, cond, &c, 3)
		aftertest(t, r, d, cond, &c, -1)
		aftertest(t, r, d, cond, &c, -2)
		aftertest(t, r, d, cond, &c, -3)
		r.Fire(d)
		for c > 6 {}
		cond.L.Lock()
		cond.Broadcast()
		cond.L.Unlock()
		for c > 0 {}
	}
}
