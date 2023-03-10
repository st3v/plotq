package queue

import (
	"errors"
	"fmt"

	"github.com/beeker1121/goque"

	v1 "github.com/st3v/plotq/api/v1"
)

type jobQueue struct {
	q *goque.Queue
}

func NewJobQueue(dir string) (*jobQueue, error) {
	q, err := goque.OpenQueue(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to open queue: %w", err)
	}

	return &jobQueue{
		q: q,
	}, nil
}

func (q *jobQueue) Enqueue(job *v1.Job) error {
	_, err := q.q.EnqueueObjectAsJSON(job)
	return err
}

func (q *jobQueue) GetAll() ([]v1.Job, error) {
	jobs := []v1.Job{}

	err := q.walkAllItems(func(item *goque.Item) error {
		job, err := jobFromItem(item)
		if err != nil {
			return err
		}

		jobs = append(jobs, *job)

		return nil
	})

	return jobs, err
}

func (q *jobQueue) Get(id string) (*v1.Job, error) {
	jobs, err := q.GetAll()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.ID == id {
			return &job, nil
		}
	}

	// Job not found, return no error
	return nil, nil
}

func (q *jobQueue) Cancel(id string) (*v1.Job, error) {
	var res *v1.Job

	err := q.walkAllItems(func(item *goque.Item) error {
		job, err := jobFromItem(item)
		if err != nil {
			return err
		}

		if job.ID != id {
			return nil
		}

		job.Status = v1.JobStatusCanceled

		if _, err := q.q.UpdateObjectAsJSON(item.ID, job); err != nil {
			return fmt.Errorf("failed to update job: %w", err)
		}

		res = job

		return stopWalk
	})

	return res, err
}

var stopWalk = errors.New("stop walk")

func (q *jobQueue) walkAllItems(callback func(item *goque.Item) error) error {
	for i := uint64(0); i < q.q.Length(); i++ {
		item, err := q.q.PeekByOffset(i)
		if err != nil {
			return fmt.Errorf("failed to peek at offset %d: %w", i, err)
		}

		if err := callback(item); err != nil {
			if err == stopWalk {
				return nil
			}
			return err
		}
	}
	return nil
}

func jobFromItem(item *goque.Item) (*v1.Job, error) {
	job := &v1.Job{}
	if err := item.ToObjectFromJSON(job); err != nil {
		return nil, fmt.Errorf("failed to decode job: %s", string(item.Value))
	}
	return job, nil
}
