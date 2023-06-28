package nonimus

import (
	"sync"
)

type Pool struct {
	Tasks   []*Task
	Workers []*Worker

	concurrency   int
	collector     chan *Task
	runBackground chan bool
	wg            sync.WaitGroup
}

func NewPool(concurrency int) *Pool {
	pool := &Pool{
		concurrency: concurrency,
		collector:   make(chan *Task, concurrency*256),
	}
	pool.run()
	return pool
}
func NewPoolCollectorSize(concurrency int, collectorSize int) *Pool {
	pool := &Pool{
		concurrency: concurrency,
		collector:   make(chan *Task, collectorSize),
	}
	pool.run()
	return pool
}

func (p *Pool) AddTask(f func()) {
	p.collector <- NewTask(f)
}

func (p *Pool) run() {
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i)
		p.Workers = append(p.Workers, worker)
		go worker.StartBackground()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
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

func AddPromise[T any](pool *Pool, executor func(resolve func(T), reject func(error))) *Promise[T] {
	if executor == nil {
		panic("executor cannot be nil")
	}

	p := &Promise[T]{
		pending: true,
		mutex:   &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}

	p.wg.Add(1)

	pool.AddTask(func() {
		defer p.handlePanic()
		executor(p.resolve, p.reject)
	})

	return p
}
