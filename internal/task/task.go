package task

import (
	"container/heap"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID         uuid.UUID
	MaxRetries int
	Priority   int
	Job        Job
	Delay      time.Duration
	RunAt      time.Time
	index      int // internal to heap implementation
}

type TaskQueue []*Task

func NewTaskQueue() TaskQueue {
	tq := make(TaskQueue, 0)
	heap.Init(&tq)

	return tq
}

func (tq *TaskQueue) AddTask(opts TaskOptions) {
	t := NewTask(opts)
	heap.Push(tq, &t)
}

func PopTask(tq *TaskQueue) (*Task, error) {
	task, ok := heap.Pop(tq).(*Task)
	if !ok {
		return nil, errors.New("failed to pop task from queue")
	}

	return task, nil
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
	Delay      time.Duration
	RunAt      time.Time
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

func (t *Task) runJob() error {
	var err error
	for range t.MaxRetries + 1 {
		err := t.Job.Execute()
		if err == nil {
			return nil
		}
	}

	return err
}

func (t *Task) GetDelay() time.Duration {
	delay := time.Duration(0)
	if t.Delay <= 0 {
		delay = 0
	}

	if t.RunAt.IsZero() {
		delay = 0
	} else {
		delay = time.Until(t.RunAt)
	}

	return delay
}

func (t *Task) Execute() error {
	delay := t.GetDelay()

	if delay > 0 {
		var err error
		time.AfterFunc(delay, func() {
			err = t.runJob()
		})

		return err
	}

	return t.runJob()
}
