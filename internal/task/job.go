package task

type JobFunc func(...any) error

type Job struct {
	job  JobFunc
	args []any
}

func NewJob(job JobFunc, args ...any) Job {
	arguments := make([]any, len(args))
	copy(arguments, args)

	return Job{
		job:  job,
		args: arguments,
	}
}

func (j *Job) Execute() error {
	return j.job(j.args...)
}
