package jobqueue

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/beeker1121/goque"

	v1 "github.com/st3v/plotq/api/v1"
)

type localQueue struct {
	q *goque.Queue // underlying queue is thread-safe
}

// localQueue implements the Queue interface.
var _ Queue = &localQueue{}

// OpenLocal opens a local queue.
func OpenLocal(dataDir string) (*localQueue, error) {
	q, err := goque.OpenQueue(filepath.Join(dataDir))
	if err != nil {
		return nil, fmt.Errorf("failed to open queue: %w", err)
	}

	return &localQueue{
		q: q,
	}, nil
}

// Close closes the queue.
func (q *localQueue) Close() error {
	return q.q.Close()
}

// Enqueue adds the given job to the queue.
func (q *localQueue) Enqueue(job *v1.Job) error {
	_, err := q.q.EnqueueObjectAsJSON(job)
	return translateGoqueError(err)
}

// GetAll returns all jobs in the queue.
func (q *localQueue) GetAll() ([]v1.Job, error) {
	jobs := []v1.Job{}

	err := q.walkAllItems(func(item *goque.Item) error {
		job, err := jobFromItem(item)
		if err != nil {
			return translateGoqueError(err)
		}

		jobs = append(jobs, *job)

		return nil
	})

	return jobs, translateGoqueError(err)
}

// Get returns the job with the given ID.
func (q *localQueue) Get(id string) (*v1.Job, error) {
	jobs, err := q.GetAll()
	if err != nil {
		return nil, translateGoqueError(err)
	}

	for _, job := range jobs {
		if job.ID == id {
			return &job, nil
		}
	}

	// Job not found, return no error
	return nil, nil
}

// Cancel marks the job with the given ID as canceled.
func (q *localQueue) Cancel(id string) (*v1.Job, error) {
	var res *v1.Job

	err := q.walkAllItems(func(item *goque.Item) error {
		job, err := jobFromItem(item)
		if err != nil {
			return translateGoqueError(err)
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

	return res, translateGoqueError(err)
}

// Peek returns the next job from the queue without removing it.
func (q *localQueue) Peek() (*v1.Job, error) {
	item, err := q.q.Peek()
	if err != nil {
		return nil, translateGoqueError(err)
	}

	return jobFromItem(item)
}

// Dequeue returns the next job from the queue.
func (q *localQueue) Dequeue() (*v1.Job, error) {
	item, err := q.q.Dequeue()
	if err != nil {
		return nil, translateGoqueError(err)
	}

	return jobFromItem(item)
}

var stopWalk = errors.New("stop walk")

func (q *localQueue) walkAllItems(callback func(item *goque.Item) error) error {
	for i := uint64(0); i < q.q.Length(); i++ {
		item, err := q.q.PeekByOffset(i)
		if err != nil {
			return err
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
		return nil, err
	}
	return job, nil
}

// translateGoqueError translates goque errors to jobqueue errors.
func translateGoqueError(err error) error {
	if errors.Is(err, goque.ErrEmpty) {
		return ErrQueueEmpty
	}
	return err
}
