# nonimus

Library for better concurrency experience in Go

## Pool (Workers)
Goroutine pool

```
pool := nonimus.NewPool(PoolSettings{
		Concurrency:       10, // 10 threads
		CollectorSize:     100, // size of buffered channel that distibutes tasks
		GoroutineStrategy: PreStartGoroutines,
		ChannelStrategy:   NativeChannelStrategy,
	})
pool.AddTask(func() {
	// do something
})
```

## Promises
Fetch some data

```
promise := NewPromise[string](func(resolve func(string), reject func(error)) {
		result, err := doSomething()
		if err != nil {
			reject(err)
			return
		}
		resolve(result)
	}).Then(func(result string) {
		// do something
	}).Catch(func(err error) {
		// do something
	})
result, err := promise.Await()
```

Add promise to pool:
AddPromise\[T any](pool *Pool, executor func(resolve func(T), reject func(error))) *Promise\[T]

## ParallelProcessArray

```
inputArray := []int{1, 2, 3, 4, 5}
outputArray := nonimus.ParallelProcessArray(nonimus.DefaultPool(), context.Background(), inputArray, func(a int, ctx context.Context) (string, error) {
	return strconv.Itoa(a * 10)
}, 2) // 2 = level of concurrency
fmt.Println(outputArray)
-> ["40", "30", "20", "10", "50"] (order can be any)

outputArray2 := nonimus.ParallelProcessArraySaveOrder(nonimus.DefaultPool(), context.Background(), inputArray, func(a int, ctx context.Context) (int, error) {
	return a * 10, nil
}, 2)
fmt.Println(outputArray2)
-> [10, 20, 30, 40, 50] (order is guaranteed to be saved)
```

## Fabric
Allows you to make a background process what will make come computations with a result

NewFabric\[T any](concurrency int, cacheSize int, scheme func() T) *Fabric\[T] {

```
helloWorldFabric := NewFabric[string](10, 1000, func() string {
		return "hello world!"
})
helloWorld1 := helloWorldFabric.Take()
helloWorld2 := helloWorldFabric.Take()
println(helloWorld2)
-> "hello world!"
println(helloWorld1)
-> "hello world!"
```
