package task

import (
	"container/heap"

	"github.com/google/uuid"
)

type Task struct {
	ID         uuid.UUID
	MaxRetries int
	Priority   int
	Job        Job
	index      int // internal to heap implementation
}

type TaskQueue []*Task

func NewTaskQueue() TaskQueue {
	tq := make(TaskQueue, 0)
	heap.Init(&tq)

	return tq
}

func (tq TaskQueue) Len() int {
	return len(tq)
}

func (tq TaskQueue) Less(i, j int) bool {
	return tq[i].Priority < tq[j].Priority
}

func (tq TaskQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
	tq[i].index, tq[j].index = i, j
}

func (tq *TaskQueue) Push(x any) {
	*tq = append(*tq, x.(*Task))
}

func (tq *TaskQueue) Pop() any {
	old := *tq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*tq = old[0 : n-1]

	return item
}

type TaskOptions struct {
	Priority   int
	MaxRetries int
	Job        Job
}

func NewTask(opts TaskOptions) Task {
	return Task{
		ID:         uuid.New(),
		MaxRetries: opts.MaxRetries,
		Priority:   opts.Priority,
		Job:        opts.Job,
	}
}

func (t *Task) Execute() error {
	var err error
	for range t.MaxRetries + 1 {
		err := t.Job.Execute()
		if err == nil {
			return nil
		}
	}

	return err
}
