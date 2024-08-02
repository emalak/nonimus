package nonimus

import (
	"sync/atomic"
	"time"
)

type Ticker struct {
	Delay     time.Duration
	Cancelled atomic.Bool
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
	ticker := &Ticker{
		Delay:     delay,
		Cancelled: atomic.Bool{},
	}
	ticker.Cancelled.Store(false)
	return ticker
}
func NewTickerFunc(delay time.Duration, f func()) *Ticker {
	ticker := NewTicker(delay)
	ticker.AddFunc(f)
	return ticker
}
