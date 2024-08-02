package nonimus

type Fabric[T any] struct {
	Cache       chan T
	Concurrency int
	CacheSize   int
}

func (f *Fabric[T]) Take() T {
	return <-f.Cache
}

func NewFabric[T any](concurrency int, cacheSize int, machine func() T) {
	cache := make(chan T, cacheSize)
	generationPool := NewPoolCollectorSize(concurrency, concurrency+1)
	generationPool.AddTask(func() {
		addTaskFunc[T](generationPool, cache, machine)
	})
}

func addTaskFunc[T any](pool *Pool, cache chan T, machine func() T) {
	for {
		if pool.settings.Concurrency == 1 {
			cache <- machine()
		} else {
			pool.AddTask(func() {
				cache <- machine()
			})
		}
	}
}
