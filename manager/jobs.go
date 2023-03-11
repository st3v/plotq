package manager

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	v1 "github.com/st3v/plotq/api/v1"
	"github.com/st3v/plotq/filestore"
	"github.com/st3v/plotq/jobqueue"
)

type jobManager struct {
	queue jobqueue.Queue
	store filestore.Store
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewJobManager(queue jobqueue.Queue, svgStore filestore.Store) *jobManager {
	return &jobManager{
		queue: queue,
		store: svgStore,
	}
}

func (m *jobManager) SubmitRequest(request *v1.JobRequest) (*v1.Job, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	request.SetDefaults()

	id := newID(request.Plotter)
	path, err := m.storeSVG(id, request)
	if err != nil {
		return nil, fmt.Errorf("failed to store SVG: %w", err)
	}

	job := &v1.Job{
		ID:          id,
		SVG:         path,
		Plotter:     request.Plotter,
		User:        request.User,
		Status:      v1.JobStatusPending,
		SubmittedAt: time.Now(),
		Settings: v1.JobSettings{
			Pagesize:    request.Pagesize,
			Velocity:    request.Velocity,
			Orientation: request.Orientation,
			Device:      request.Device,
		},
	}

	if err := m.queue.Enqueue(job); err != nil {
		return nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	return job, nil
}

func (m *jobManager) GetJobs() ([]v1.Job, error) {
	return m.queue.GetAll()
}

func (m *jobManager) GetJob(id string) (*v1.Job, error) {
	return m.queue.Get(id)
}

func (m *jobManager) DeleteJob(id string) (*v1.Job, error) {
	return m.queue.Cancel(id)
}

func (m *jobManager) storeSVG(id string, request *v1.JobRequest) (string, error) {
	svg, err := request.SVG.Open()
	if err != nil {
		return "", fmt.Errorf("could not open file %s: %w", request.SVG.Filename, err)
	}
	defer svg.Close()

	path := fmt.Sprintf("%s.svg", id)

	written, err := m.store.Put(path, svg)
	if err != nil {
		return "", err
	}

	if written != request.SVG.Size {
		return "", errors.New("size mismatch")
	}

	return path, nil
}

const alphanumeric = "0123456789abcdefghijklmnopqrstuvwxyz"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, randString(8))
}
