package nonimus

import (
	"encoding/json"
	"errors"
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

// TODO
func all(promises ...Promise[any]) []any {
	return []any{}
}

func AddPromise[T any](pool *Pool, executor func(resolve func(T), reject func(error))) *Promise[T] {
	if executor == nil {
		// log("Executor cannot be nil!")
		return &Promise[T]{
			pending: false,
			mutex:   nil,
			wg:      nil,
		}
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

func NewPromise[T any](executor func(resolve func(T), reject func(error))) *Promise[T] {
	return AddPromise(defaultPool, executor)
}
func NewPromiseFunc[T any, Argument any](executor func(args Argument) (T, error), args Argument) *Promise[T] {
	if executor == nil {
		// log("Executor cannot be nil!")
		return &Promise[T]{
			pending: false,
			mutex:   nil,
			wg:      nil,
		}
	}

	p := &Promise[T]{
		pending: true,
		mutex:   &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}

	p.wg.Add(1)

	defaultPool.AddTask(func() {
		defer p.handlePanic()
		result, err := executor(args)
		if err != nil {
			p.reject(err)
			return
		}
		p.resolve(result)
	})

	return p
}
func NewPromiseFuncReturnJSON[T any, Argument any](executor func(args Argument) (T, error), args Argument) *Promise[string] {
	if executor == nil {
		// log("Executor cannot be nil!")
		return &Promise[string]{
			pending: false,
			mutex:   nil,
			wg:      nil,
		}
	}

	p := &Promise[string]{
		pending: true,
		mutex:   &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}

	p.wg.Add(1)

	defaultPool.AddTask(func() {
		defer p.handlePanic()
		result, err := executor(args)
		if err != nil {
			p.reject(err)
			return
		}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			p.reject(err)
			return
		}
		p.resolve(string(jsonResult))
	})

	return p
}

// Await blocks until the promise is resolved or rejected.
func (p *Promise[T]) Await() (T, error) {
	if p.wg == nil {
		return p.result, errors.New("Executor is null! ")
	}
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
	if err == nil {
		return
	}
	if validErr, ok := err.(error); ok {
		// log(validErr.Error())
		p.reject(fmt.Errorf("panic recovery: %w", validErr))
	} else {
		// log("panic recovery: %v\n", err)
		p.reject(fmt.Errorf("panic recovery: %v", err))
	}
}

func (p *Promise[T]) resolve(resolution T) {
	p.mutex.Lock()
	if !p.pending {
		p.mutex.Unlock()
		return
	}

	p.result = resolution
	p.pending = false
	if p.then != nil {
		p.then(p.result)
	}

	p.mutex.Unlock()
	p.wg.Done()
}

func (p *Promise[T]) reject(err error) {
	p.mutex.Lock()
	if !p.pending {
		p.mutex.Unlock()
		return
	}

	p.err = err
	p.pending = false
	if p.catch != nil {
		p.catch(p.err)
	}

	p.mutex.Unlock()
	p.wg.Done()
}
