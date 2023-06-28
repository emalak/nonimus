# nonimus

Library for better concurrency experience in Go

## Pool (Workers)
Goroutine pool

```
pool := nonimus.NewPool(10) // 10 threads
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
println(helloWorld1)
```
