package task

import (
	"fmt"
	"reflect"
)

type Job struct {
	job  any
	args []any
}

func NewJob(job any, args ...any) Job {
	arguments := make([]any, len(args))
	copy(arguments, args)

	return Job{
		job:  job,
		args: arguments,
	}
}

func (j *Job) Execute() error {
	jobValue := reflect.ValueOf(j.job)
	if jobValue.Kind() != reflect.Func {
		return fmt.Errorf("job must be a function")
	}

	var in []reflect.Value
	for _, arg := range j.args {
		in = append(in, reflect.ValueOf(arg))
	}

	// check if any error is returned, used to determine if task should be retried
	results := jobValue.Call(in)
	if len(results) > 0 && results[0].CanInterface() {
		if err, ok := results[0].Interface().(error); ok {
			return err
		}
	}
	return nil
}
