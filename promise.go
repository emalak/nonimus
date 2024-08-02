package nonimus

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Promise[T any] struct {
	result T
	err    error

	then  func(t T)
	catch func(err error)

	mutex sync.Mutex
}

// TODO
func All(promises ...*Promise[any]) {

}

func AddPromise[T any](pool *Pool, executor func() (T, error)) *Promise[T] {
	if executor == nil {
		// fmt.Println("AddPromise: Executor cannot be nil!")
		return &Promise[T]{
			mutex: sync.Mutex{},
		}
	}

	p := &Promise[T]{
		mutex: sync.Mutex{},
	}
	p.mutex.Lock()

	pool.AddTask(func() {
		defer p.deferFunc()
		t, err := executor()
		if err != nil {
			p.reject(err)
			return
		}
		p.resolve(t)
	})

	return p
}

func NewPromise[T any](executor func() (T, error)) *Promise[T] {
	return AddPromise(DefaultPool(), executor)
}
func NewPromiseFunc[T any, Arg any](executor func(args Arg) (T, error), args Arg) *Promise[T] {
	return NewPromise(func() (T, error) {
		return executor(args)
	})
}
func NewPromiseFuncReturnJSON[T any, Argument any](executor func(args Argument) (T, error), args Argument) *Promise[string] {
	return NewPromise(func() (string, error) {
		t, err := executor(args)
		if err != nil {
			return "", err
		}
		jsonBytes, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	})
}

// Await blocks until the promise is resolved or rejected.
func (p *Promise[T]) Await() (T, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.result, p.err
}

// Then executes function if no err
func (p *Promise[T]) Then(f func(t T)) *Promise[T] {
	locked := p.mutex.TryLock()
	p.then = f
	if locked {
		defer p.mutex.Unlock()
		if p.err == nil {
			p.then(p.result)
		}
	}
	return p
}

// Catch executes if err != nil
func (p *Promise[T]) Catch(f func(err error)) *Promise[T] {
	locked := p.mutex.TryLock()
	p.catch = f
	if locked {
		defer p.mutex.Unlock()
		if p.err != nil {
			p.catch(p.err)
		}
	}
	return p
}

func (p *Promise[T]) deferFunc() {
	p.handlePanic()
	p.mutex.Unlock()
}

func (p *Promise[T]) handlePanic() {
	err := recover()
	if err == nil {
		return
	}
	if validErr, ok := err.(error); ok {
		fmt.Println(validErr.Error())
		p.reject(fmt.Errorf("promise panic recovery: %w", validErr))
	} else {
		fmt.Println("promise panic recovery: %v\n", err)
		p.reject(fmt.Errorf("promise panic recovery: %v", err))
	}
}

func (p *Promise[T]) resolve(resolution T) {
	p.result = resolution
	if p.then != nil {
		p.then(p.result)
	}
}

func (p *Promise[T]) reject(err error) {
	p.err = err
	if p.catch != nil {
		p.catch(p.err)
	}
}
