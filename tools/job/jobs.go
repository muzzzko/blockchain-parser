package job

import (
	"context"
	"sync"
)

type Jobs []*Job

func (jobs *Jobs) Add(job *Job) {
	*jobs = append(*jobs, job)
}

func (jobs *Jobs) Start(ctx context.Context) {
	for _, job := range *jobs {
		job.Start(ctx)
	}
}

func (jobs *Jobs) Stop() {
	wg := sync.WaitGroup{}

	for _, job := range *jobs {
		job := job

		wg.Add(1)
		go func() {
			defer wg.Done()
			job.Stop()
		}()
	}

	wg.Wait()
}
