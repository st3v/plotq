package jobqueue

import v1 "github.com/st3v/plotq/api/v1"

type Queue interface {
	Enqueue(job *v1.Job) error
	GetAll() ([]v1.Job, error)
	Get(id string) (*v1.Job, error)
	Cancel(id string) (*v1.Job, error)
}
