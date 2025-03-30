package main

import (
	"fmt"

	"github.com/jackparsonss/chronos/internal/task"
	"github.com/jackparsonss/chronos/pkg/scheduler"
)

func main() {
	s := scheduler.NewScheduler()

	s.ScheduleTask(scheduler.TaskOpts{
		MaxRetries: 3,
		Priority:   1,
		Job: task.NewJob(func(x, y int) error {
			fmt.Println(x + y)

			return nil
		}, 10, 20),
	})

	go s.Run()

	select {}
}
