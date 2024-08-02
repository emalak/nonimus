package nonimus

import (
	"context"
	"database/sql"
	"sync"
)

type Parallex[MaterialType, ResultType any] struct { // TODO
	result chan ResultType
}

type ResultAndErrorStruct[R any] struct {
	Result R
	Error  error
}

type ResultSlice[T any] []T

func (slice *ResultSlice[T]) AppendResultOrReturnError(result T, err error) error {
	if err != nil {
		*slice = append(*slice, result)
		return nil
	} else {
		return err
	}
}
func Parallel[M any](pool *Pool, materials []M, process func(m M, ctx context.Context, cancelFunc context.CancelFunc, mx *sync.Mutex) error, concurrency int) error {
	materialsLen := len(materials)
	if materialsLen == 0 {
		return nil
	}
	if concurrency > materialsLen {
		concurrency = materialsLen
	}
	cancelContext, cancelFunc := context.WithCancel(context.Background())
	materialsChan := make(chan M, materialsLen)
	defer func() {
		cancelFunc()
		close(materialsChan)
	}()

	for i := 0; i < materialsLen; i++ {
		materialsChan <- materials[i]
	}
	var wg sync.WaitGroup
	var err error

	wg.Add(concurrency)
	mx := &sync.Mutex{}

	for i := 0; i < concurrency; i++ {
		pool.AddTask(func() {
			defer wg.Done()
			var m M
			var ok bool
			var localErr error
			for {
				select {
				case <-cancelContext.Done():
					return
				case m, ok = <-materialsChan:
					if !ok {
						return
					}
					localErr = process(m, cancelContext, cancelFunc, mx)
					if localErr != nil {
						mx.Lock()
						cancelFunc()
						err = localErr
						mx.Unlock()
						return
					}
					continue
				default:
					return
				}
			}
		})
	}

	wg.Wait()
	return err
}

func ParallelProcessArray[M any, R any](pool *Pool, ctx context.Context, materials []M, process func(m M, ctx context.Context) (R, error), concurrency int) ([]R, error) {
	materialsLen := len(materials)
	if materialsLen == 0 {
		return []R{}, nil
	}
	if concurrency < 2 {
		var err error
		result := make([]R, materialsLen)
		for i := 0; i < materialsLen; i++ {
			result[i], err = process(materials[i], ctx)
			if err != nil {
				return result, err
			}
		}
		return result, nil
	}
	if concurrency > materialsLen {
		concurrency = materialsLen
	}
	materialsChan := make(chan M, materialsLen)
	cancelContext, cancelFunc := context.WithCancel(ctx)
	defer func() {
		cancelFunc()
		close(materialsChan)
	}()

	for i := 0; i < materialsLen; i++ {
		materialsChan <- materials[i]
	}

	result := make([]R, materialsLen)
	var resultIndex int
	var wg sync.WaitGroup
	var err error
	mx := sync.Mutex{}

	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		pool.AddTask(func() {
			defer wg.Done()
			var r R
			var localErr error
			var m M
			var ok bool
			for {
				select {
				case <-cancelContext.Done():
					return
				case m, ok = <-materialsChan:
					if !ok {
						return
					}
					r, localErr = process(m, cancelContext)
					if localErr != nil {
						mx.Lock()
						err = localErr
						cancelFunc()
						mx.Unlock()
						return
					}
					mx.Lock()
					result[resultIndex] = r
					resultIndex++
					mx.Unlock()
					continue
				default:
					return
				}
			}
		})
	}

	wg.Wait()
	return result, err
}

type OrderedMaterial[M any] struct {
	Item  M
	Index int
}
type OrderedResultAndErrorStruct[R any] struct {
	Result R
	Error  error
	Index  int
}

func ParallelProcessArraySaveOrder[M any, R any](pool *Pool, ctx context.Context, materials []M, process func(m M, ctx context.Context) (R, error), concurrency int) ([]R, error) {
	materialsLen := len(materials)
	if materialsLen == 0 {
		return []R{}, nil
	}
	if concurrency < 2 {
		var err error
		result := make([]R, materialsLen)
		for i := 0; i < materialsLen; i++ {
			result[i], err = process(materials[i], ctx)
			if err != nil {
				return result, err
			}
		}
		return result, nil
	}
	if concurrency > materialsLen {
		concurrency = materialsLen
	}
	materialsChan := make(chan OrderedMaterial[M], materialsLen)
	cancelContext, cancelFunc := context.WithCancel(ctx)
	defer func() {
		cancelFunc()
		close(materialsChan)
	}()

	for i := 0; i < materialsLen; i++ {
		materialsChan <- OrderedMaterial[M]{Item: materials[i], Index: i}
	}
	result := make([]R, materialsLen)
	var wg sync.WaitGroup
	var err error
	mx := sync.Mutex{}

	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		pool.AddTask(func() {
			defer wg.Done()
			var r R
			var localErr error
			var m OrderedMaterial[M]
			var ok bool
			for {
				select {
				case <-cancelContext.Done():
					return
				case m, ok = <-materialsChan:
					if !ok {
						return
					}
					r, localErr = process(m.Item, cancelContext)
					if localErr != nil {
						mx.Lock()
						err = localErr
						cancelFunc()
						mx.Unlock()
						return
					}
					mx.Lock()
					result[m.Index] = r
					mx.Unlock()
					continue
				default:
					return
				}
			}
		})
	}

	wg.Wait()
	return result, err
}

// TODO
func todoParallelProcessRows[R any](pool *Pool, rows *sql.Rows, process func() (R, error), concurrency int) ([]R, error) {
	return nil, nil
}
