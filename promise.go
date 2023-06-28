package nonimus

import (
	"fmt"
	"sync"
)

type Promise[T any] struct {
	result T
	err    error

	then  func(t T)
	catch func(err error)

	pending bool
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

func NewPromise[T any](executor func(resolve func(T), reject func(error))) *Promise[T] {
	if executor == nil {
		panic("executor cannot be nil")
	}

	p := &Promise[T]{
		pending: true,
		mutex:   &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}

	p.wg.Add(1)

	go func() {
		defer p.handlePanic()
		executor(p.resolve, p.reject)

	}()

	return p
}

// Await blocks until the promise is resolved or rejected.
func (p *Promise[T]) Await() (T, error) {
	p.wg.Wait()
	return p.result, p.err
}

// Then executes function if no err
func (p *Promise[T]) Then(f func(t T)) *Promise[T] {
	p.then = f
	return p
}

// Catch executes if err != nil
func (p *Promise[T]) Catch(f func(err error)) *Promise[T] {
	p.catch = f
	return p
}

func (p *Promise[T]) handlePanic() {
	err := recover()
	if validErr, ok := err.(error); ok {
		p.reject(fmt.Errorf("panic recovery: %w", validErr))
	} else {
		p.reject(fmt.Errorf("panic recovery: %+v", err))
	}
}

func (p *Promise[T]) resolve(resolution T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.pending {
		return
	}

	p.result = resolution
	p.pending = false
	p.then(p.result)

	p.wg.Done()
}

func (p *Promise[T]) reject(err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.pending {
		return
	}

	p.err = err
	p.pending = false
	p.catch(p.err)

	p.wg.Done()
}
