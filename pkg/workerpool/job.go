package workerpool

import "context"

type JobDescriptor struct {
	ID int
}

type ExecutionFn func(context.Context, interface{}) (interface{}, error)

type Job struct {
	Descriptor JobDescriptor
	ExecFn     func(context.Context, interface{}) (interface{}, error)
	Args       interface{}
}

type Result struct {
	Value      interface{}
	Descriptor JobDescriptor
	Err        error
}

func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Value:      value,
		Descriptor: j.Descriptor,
	}
}
