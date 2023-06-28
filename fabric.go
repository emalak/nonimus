package nonimus

type Fabric[T any] struct {
	Cache       chan T
	Concurrency int
	CacheSize   int
}

func (f *Fabric[T]) Take() T {
	return <-f.Cache
}

func NewFabric[T any](concurrency int, cacheSize int, scheme func() T) *Fabric[T] {
	cache := make(chan T, cacheSize)
	generationPool := NewPoolCollectorSize(concurrency, concurrency*16)
	generationPool.AddTask(func() {
		addTaskFunc[T](generationPool, cache, scheme)
	})
	return &Fabric[T]{Cache: cache, Concurrency: concurrency, CacheSize: cacheSize}
}

func addTaskFunc[T any](pool *Pool, cache chan T, scheme func() T) {
	for {
		if pool.concurrency == 1 {
			cache <- scheme()
		} else {
			pool.AddTask(func() {
				cache <- scheme()
			})
		}
	}
}
