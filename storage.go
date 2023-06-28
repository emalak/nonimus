package nonimus

import (
	"sync"
	"time"
)

type Storage[T any] interface {
	Get(key interface{}) (value *T)
	Remove(key interface{}) (value *T)
	Set(key interface{}, value T)
	Reset()
}

type TimedStorage[T any] interface {
	Get(key interface{}) (value *T)
	Remove(key interface{}) (value *T)
	Set(key interface{}, value T, lifetime time.Duration)
	Reset()
}

type MemoryStorage[T any] struct {
	driver map[string]T
	mx     sync.RWMutex
}

func (s *MemoryStorage[T]) Get(key string) *T {
	s.mx.RLock()
	r, ok := s.driver[key]
	s.mx.RUnlock()
	if !ok {
		return nil
	}
	return &r
}
func (s *MemoryStorage[T]) Remove(key string) *T {
	s.mx.RLock()
	r, ok := s.driver[key]
	s.mx.RUnlock()
	if !ok {
		return nil
	}
	s.mx.Lock()
	delete(s.driver, key)
	s.mx.Unlock()
	return &r
}
func (s *MemoryStorage[T]) Set(key string, value T) {
	s.mx.Lock()
	s.driver[key] = value
	s.mx.Unlock()
}
func (s *MemoryStorage[T]) Reset() {
	s.mx.Lock()
	for key, _ := range s.driver {
		delete(s.driver, key)
	}
	s.mx.Unlock()
}
