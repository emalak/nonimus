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
func (shard *Shard) AsyncGetString(key Key, column string) *nonimus.Promise[string] {
	return nonimus.AddPromise(pool, func(resolve func(GetResult[string]), reject func(error)) {
		var result string
		err, found := shard.Get(Keys{key}, SelectColumns{{column, &result}})
		if err != nil {
			reject(err)
			return
		}
		resolve(result)
	})
}
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
