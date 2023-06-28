package nonimus

import (
	"sync/atomic"
	"time"
)

type Ticker struct {
	Delay     time.Duration
	Cancelled *atomic.Bool
}

func (t *Ticker) Cancel() {
	t.Cancelled.Store(true)
}

func (t *Ticker) AddFunc(f func()) {
	go func() {
		for range time.Tick(t.Delay) {
			if t.Cancelled.Load() {
				break
			}
			f()
		}
	}()
}

func NewTicker(delay time.Duration) *Ticker {
	var cancelled *atomic.Bool
	cancelled.Store(false)
	return &Ticker{Delay: delay, Cancelled: cancelled}
}
func NewTickerFunc(delay time.Duration, f func()) *Ticker {
	var cancelled *atomic.Bool
	cancelled.Store(false)
	ticker := &Ticker{Delay: delay, Cancelled: cancelled}
	ticker.AddFunc(f)
	return ticker
}
