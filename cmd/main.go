package main

import (
	"fmt"
	"time"

	"github.com/jackparsonss/chronos/internal/task"
	"github.com/jackparsonss/chronos/pkg/scheduler"
)

func main() {
	s := scheduler.NewScheduler()

	s.ScheduleTask(scheduler.TaskOpts{
		MaxRetries: 3,
		Priority:   1,
		Delay:      time.Second * 2,
		Job: task.NewJob(func(x, y int) {
			fmt.Println(x + y)
		}, 10, 20),
	})

	go s.Run()

	select {}
}
