package nonimus

var defaultPool *Pool

func init() {
	defaultPool = NewPool(PoolSettings{
		Concurrency:       128,
		CollectorSize:     32000,
		GoroutineStrategy: PreStartGoroutines,
		ChannelStrategy:   NativeChannelStrategy,
	}) // TODO config
}

func DefaultPool() *Pool {
	return defaultPool
}

func Recover(additional ...string) {
	recover()
	/*err := recover()
	if realErr, ok := err.(error); ok {
		TODO add custom logging of panics
	}*/
}
