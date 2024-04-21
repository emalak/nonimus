package nonimus

import (
	"github.com/alphadose/zenq/v2"
	"sync"
)

type Task func()

type Pool struct {
	Workers []*Worker

	settings      PoolSettings
	runBackground chan bool
	wg            sync.WaitGroup

	// Channel
	collector chan Task
	zenq      *zenq.ZenQ[Task] // TODO for LowestLatencyChannel
}

type PoolGoroutineStrategy int

const (
	PreStartGoroutines PoolGoroutineStrategy = iota
)

type PoolChannelStrategy int

const (
	NativeChannelStrategy PoolChannelStrategy = iota

	// LowestLatencyChannel is not implemented yet
	LowestLatencyChannelStrategy
)

type PoolSettings struct {
	Concurrency   int
	CollectorSize int

	GoroutineStrategy PoolGoroutineStrategy
	ChannelStrategy   PoolChannelStrategy
}

func NewPool(settings PoolSettings) *Pool {
	pool := &Pool{settings: settings}
	switch pool.settings.ChannelStrategy {
	case NativeChannelStrategy:
		{
			pool.collector = make(chan Task, settings.CollectorSize)
		}
	case LowestLatencyChannelStrategy:
		{
			pool.zenq = zenq.New[Task](10)
		}
	}
	pool.run()
	return pool
}
func NewPoolCollectorSize(concurrency int, collectorSize int) *Pool {
	return NewPool(PoolSettings{
		Concurrency:       concurrency,
		CollectorSize:     collectorSize,
		GoroutineStrategy: PreStartGoroutines,
		ChannelStrategy:   NativeChannelStrategy,
	})
}

func (p *Pool) AddTask(f func()) {
	switch p.settings.ChannelStrategy {
	case NativeChannelStrategy:
		{
			p.collector <- f
			return
		}
	case LowestLatencyChannelStrategy:
		{
			p.zenq.Write(f)
			return
		}
	default:
		{
			p.collector <- f
			return
		}
	}
}

func (p *Pool) run() {
	for i := 1; i <= p.settings.Concurrency; i++ {
		worker := NewWorker(p.collector, i)
		p.Workers = append(p.Workers, worker)
		go worker.StartBackground()
	}

	p.runBackground = make(chan bool)
}

func (p *Pool) Wait() {
	<-p.runBackground
}

func (p *Pool) Stop() {
	for i := range p.Workers {
		p.Workers[i].Stop()
	}
	p.runBackground <- true
}
