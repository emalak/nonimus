package nonimus

import (
	"context"
	"errors"
)

type ResultAndErrorStruct[R any] struct {
	Result R
	Error  error
}

func ResultAndError[R any](r R, err error) ResultAndErrorStruct[R] {
	return ResultAndErrorStruct[R]{Result: r, Error: err}
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

func ParallelProcessArray[M any, R any](pool *Pool, ctx context.Context, materials []M, process func(m M, ctx context.Context) (R, error), concurrency int) ([]R, error) {
	materialsLen := len(materials)
	if materialsLen == 0 {
		return []R{}, nil
	}
	if concurrency > materialsLen {
		concurrency = materialsLen
	}
	materialsChan := make(chan M, materialsLen)
	resultChan := make(chan ResultAndErrorStruct[R])
	cancelContext, cancelFunc := context.WithCancel(ctx)
	defer func() {
		cancelFunc()
		close(materialsChan)
	}()

	for i := 0; i < concurrency; i++ {
		materialsChan <- materials[i]
		pool.AddTask(func() {
			var m M
			for {
				select {
				case <-cancelContext.Done():
					return
				case m = <-materialsChan:
					resultChan <- ResultAndError(process(m, cancelContext))
					continue
				}
			}
		})
	}
	for i := concurrency; i < materialsLen; i++ {
		materialsChan <- materials[i]
	}

	result := make([]R, materialsLen)
	var r ResultAndErrorStruct[R]
	for i := 0; i < materialsLen; i++ {
		select {
		case <-cancelContext.Done():
			{
				return result, errors.New("Timeout ")
			}
		case r = <-resultChan:
			{
				if r.Error != nil {
					// TODO add handling error function
					return result, r.Error
				}
				// TODO add handling nil function
				result[i] = r.Result
			}
		}
	}
	return result, nil
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
	if concurrency > materialsLen {
		concurrency = materialsLen
	}
	materialsChan := make(chan OrderedMaterial[M], materialsLen)
	resultChan := make(chan OrderedResultAndErrorStruct[R])
	cancelContext, cancelFunc := context.WithCancel(ctx)
	defer func() {
		cancelFunc()
		close(materialsChan)
	}()

	for i := 0; i < concurrency; i++ {
		materialsChan <- OrderedMaterial[M]{Item: materials[i], Index: i}
		pool.AddTask(func() {
			var m OrderedMaterial[M]
			var r R
			var err error
			for {
				select {
				case <-cancelContext.Done():
					return
				case m = <-materialsChan:
					r, err = process(m.Item, cancelContext)
					resultChan <- OrderedResultAndErrorStruct[R]{Index: m.Index, Result: r, Error: err}
					continue
				}
			}
		})
	}
	for i := concurrency; i < materialsLen; i++ {
		materialsChan <- OrderedMaterial[M]{Item: materials[i], Index: i}
	}

	result := make([]R, materialsLen)
	var r OrderedResultAndErrorStruct[R]
	for i := 0; i < materialsLen; i++ {
		select {
		case <-cancelContext.Done():
			{
				return result, errors.New("Timeout ")
			}
		case r = <-resultChan:
			{
				if r.Error != nil {
					return result, r.Error
				}
				result[r.Index] = r.Result
			}
		}
	}
	return result, nil
}

func undoParallel(pool *Pool, process func() error, concurrency int) error {
	return nil
}
