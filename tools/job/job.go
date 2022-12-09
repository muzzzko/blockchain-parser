package job

import (
	"context"
	"log"
	"runtime/debug"
	"time"
)

type Job struct {
	run      func(ctx context.Context) error
	name     string
	interval time.Duration

	stop chan struct{}
}

func NewJob(
	run func(ctx context.Context) error,
	name string,
	interval time.Duration,
) *Job {
	return &Job{
		run:      run,
		name:     name,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (j *Job) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)

	go func() {
		log.Printf("job %s started", j.name)

		for {
			select {
			case <-j.stop:
				break
			default:
			}

			select {
			case <-ticker.C:
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Println("stacktrace from panic: \n" + string(debug.Stack()))
						}

						log.Printf("run %s", j.name)

						if err := j.run(ctx); err != nil {
							log.Printf("fail run job %s: %s", j.name, err)
						}

						log.Printf("finish %s", j.name)
					}()
				}()
			case <-j.stop:
				break
			}
		}
	}()
}

func (j *Job) Stop() {
	j.stop <- struct{}{}
	close(j.stop)

	log.Printf("job %s stopped", j.name)
}
