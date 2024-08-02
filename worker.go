package nonimus

type Worker struct {
	ID       int
	taskChan chan Task
}

func NewWorker(channel chan Task, ID int) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: channel,
	}
}

func (w *Worker) StartBackground() {
	for {
		select {
		case task := <-w.taskChan:
			(func() {
				defer Recover()
				task()
			})()
			//case <-w.ctx.Done():
			//	return
		}
	}
}

func (w *Worker) Stop() {
	// TODO
}
