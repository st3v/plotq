package worker

import (
	"context"
	"log"

	v1 "github.com/st3v/plotq/api/v1"
)

type Spooler interface {
	Process(job v1.Job) (sent int64, err error)
	Incoming(ctx context.Context) <-chan v1.Job
}

// Run runs a worker loop.
func Run(ctx context.Context, spooler Spooler) error {
	jobs := spooler.Incoming(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case job := <-jobs:
			log.Printf("processing job %s...", job.ID)
			job.Status = v1.JobStatusProcessing
			sent, err := spooler.Process(job)
			if err != nil {
				log.Printf("job %s failed: %v", job.ID, err)
				job.Error = err.Error()
				job.Status = v1.JobStatusFailed
			} else {
				log.Printf("job %s succeeded: %d bytes sent to plotter", job.ID, sent)
				job.Status = v1.JobStatusSucceeded
			}

			// todo: update job
		}
	}
}
