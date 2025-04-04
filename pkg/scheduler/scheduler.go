package scheduler

import (
	"github.com/jackparsonss/chronos/internal/task"
)

type Scheduler struct {
	Tasks task.TaskQueue
}

func NewScheduler() Scheduler {
	return Scheduler{
		Tasks: task.NewTaskQueue(),
	}
}

type TaskOpts task.TaskOptions

func (s *Scheduler) ScheduleTask(opts TaskOpts) {
	s.Tasks.AddTask(task.TaskOptions{
		Priority:   opts.Priority,
		MaxRetries: opts.MaxRetries,
		Job:        opts.Job,
		Delay:      opts.Delay,
		RunAt:      opts.RunAt,
	})
}

func (s *Scheduler) ExecuteNextTask() error {
	task, err := task.PopTask(&s.Tasks)
	if err != nil {
		return err
	}

	go task.Execute()

	return nil
}

func (s *Scheduler) HasTasks() bool {
	return s.Tasks.Len() > 0
}

func (s *Scheduler) Run() {
	for {
		if s.HasTasks() {
			s.ExecuteNextTask()
		}
	}
}
