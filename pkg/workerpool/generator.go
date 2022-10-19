package workerpool

func (wp *WorkerPool) GenerateFrom(jobsBulk []Job) {
	for i := range jobsBulk {
		wp.Working++
		wp.jobs <- jobsBulk[i]
	}
}

func (wp *WorkerPool) GenerateFromJob(job Job) {
	wp.Working++
	wp.jobs <- job
}
