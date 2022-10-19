package workerpool

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

var dummyErr error = errors.New("dummy error")

const numberOfJobs int = 2

func dummyExecFn(ctx context.Context, job interface{}) (interface{}, error) {
	if job.(int)%2 == 0 {
		return nil, dummyErr
	}
	return job.(int) * 2, nil
}

func testJobs() []Job {
	jobs := make([]Job, 0, 10)
	for i := 0; i < numberOfJobs; i++ {
		jobs = append(jobs, Job{Descriptor: JobDescriptor{ID: i}, ExecFn: dummyExecFn, Args: i})
	}
	return jobs
}

func TestGenerateFrom(t *testing.T) {
	wp := New(10)
	wp.GenerateFrom(testJobs())
	for i := 0; i < numberOfJobs; i++ {
		job := <-wp.jobs
		if reflect.ValueOf(job.ExecFn).Pointer() != reflect.ValueOf(dummyExecFn).Pointer() {
			t.Errorf("Expected job execution function")
		}

		if job.Args != job.Descriptor.ID {
			t.Errorf("Expected job argument %d, got %d", job.Descriptor.ID, job.Args)
		}
	}
}

func TestGenerateFromJob(t *testing.T) {
	wp := New(10)
	for i := 0; i < numberOfJobs; i++ {
		wp.GenerateFromJob(Job{Descriptor: JobDescriptor{ID: i}, ExecFn: dummyExecFn, Args: i})
	}
	for i := 0; i < numberOfJobs; i++ {
		job := <-wp.jobs
		if reflect.ValueOf(job.ExecFn).Pointer() != reflect.ValueOf(dummyExecFn).Pointer() {
			t.Errorf("Expected job execution function")
		}

		if job.Args != job.Descriptor.ID {
			t.Errorf("Expected job argument %d, got %d", job.Descriptor.ID, job.Args)
		}
	}
}

func TestWorkerPool_Run(t *testing.T) {
	wp := New(10)
	wp.GenerateFrom(testJobs())

	go wp.Run(context.TODO())

	var tests = map[int]struct {
		jobID int
		err   error
		want  interface{}
	}{
		0: {0, dummyErr, nil},
		1: {1, nil, 2},
	}

	for i := 0; i < numberOfJobs; i++ {
		result := <-wp.results

		if result.Err != tests[result.Descriptor.ID].err {
			t.Errorf("Expected result error %v, got %v", tests[result.Descriptor.ID].err, result.Err)
		}

		if result.Value != tests[result.Descriptor.ID].want {
			t.Errorf("Expected result value %d, got %d", tests[result.Descriptor.ID].want, result.Value)
		}
	}
}

func TestWorkerPoolNew(t *testing.T) {
	wp := New(10)
	if wp.workersCount != 10 {
		t.Errorf("Expected 10 workers, got %d", wp.workersCount)
	}
}

func TestWorkerPool_RunDoneContext(t *testing.T) {
	wp := New(10)
	wp.GenerateFrom(testJobs())
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	cancelCtx()

	go wp.Run(ctx)

	result := <-wp.results

	if result.Err == nil {
		t.Errorf("Expected result error %v, got nil", context.Canceled)
	}

}

func TestWorkerPool_RunNoResults(t *testing.T) {
	wp := New(10)
	t.Log("Starting workers")
	ctx := context.Background()
	wp.CloseJobsChannel()

	go wp.Run(ctx)

	var tests = map[int]struct {
		jobID int
		err   error
		want  interface{}
	}{
		0: {0, nil, nil},
		1: {1, nil, 2},
	}

	for i := 0; i < numberOfJobs; i++ {
		result := <-wp.results

		if result.Err != tests[result.Descriptor.ID].err {
			t.Errorf("Expected result error %v, got %v", tests[result.Descriptor.ID].err, result.Err)
		}

		if result.Value != tests[result.Descriptor.ID].want {
			t.Errorf("Expected result value %d, got %d", tests[result.Descriptor.ID].want, result.Value)
		}
	}
}

func TestWorkerPool_Results(t *testing.T) {
	wp := New(10)

	wp.GenerateFrom(testJobs())
	ctx := context.Background()
	go wp.Run(ctx)

	results := wp.Results()
	if results == nil {
		t.Errorf("Expected results channel, got nil")
	}
}

func TestWorkerPool_CloseWorkChannel(t *testing.T) {
	wp := New(10)
	wp.CloseJobsChannel()
	select {
	case <-wp.jobs:
	default:
		t.Errorf("Expected closed jobs channel, got open channel")
	}
}

func TestWorkerPool_GetQueueSize(t *testing.T) {
	wp := New(10)
	wp.GenerateFrom(testJobs())
	ctx := context.Background()
	go wp.Run(ctx)

	size := wp.GetQueueSize()
	if size != numberOfJobs {
		t.Errorf("Expected queue size %d, got %d", numberOfJobs, size)
	}
}
