package nonimus

import "fmt"

var defaultPool *Pool

type Settings struct {
	DefaultPoolSettings *PoolSettings
}

func Init(s Settings) {
	if s.DefaultPoolSettings == nil {
		s.DefaultPoolSettings = &PoolSettings{
			Name:          "DEFAULT_POOL",
			Concurrency:   4,
			CollectorSize: 16,
		}
	}
	defaultPool = NewPool(*s.DefaultPoolSettings)
}

func DefaultPool() *Pool {
	return defaultPool
}

func Recover(additional ...string) {
	err := recover()
	if realErr, ok := err.(error); ok {
		fmt.Println("PANIC! ", realErr, additional)
	}
}
