package nonimus

type Worker struct {
	ID       int
	taskChan chan Task
	quit     chan bool
}

func NewWorker(channel chan Task, ID int) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
	}
}

func (wr *Worker) StartBackground() {
	for {
		select {
		case task := <-wr.taskChan:
			(func() {
				defer Recover()
				task()
			})()
		case <-wr.quit:
			return
		}
	}
}

func (wr *Worker) Stop() {
	go func() {
		wr.quit <- true
	}()
}
