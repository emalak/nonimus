package nonimus

type Task struct {
	f func()
}

func NewTask(f func()) *Task {
	return &Task{f: f}
}

func process(workerID int, task *Task) {
	task.f()
}
