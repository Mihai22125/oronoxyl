package workerpool

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Working      int
	Done         chan struct{}
}

func New(wcount int) WorkerPool {
	return WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, 1000),
		results:      make(chan Result, 1000),
		Done:         make(chan struct{}),
	}
}

func (wp WorkerPool) Results() chan Result {
	return wp.results
}

func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp WorkerPool) CloseJobsChannel() {
	select {
	case <-wp.jobs:
	default:
		close(wp.jobs)
	}
}

func (wp WorkerPool) GetQueueSize() int {
	return len(wp.jobs)
}
