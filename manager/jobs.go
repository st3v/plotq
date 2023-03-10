package manager

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"

	v1 "github.com/st3v/plotq/api/v1"
	"github.com/st3v/plotq/jobqueue"
)

var uploadDir = path.Join("data", "uploads")

type jobManager struct {
	queue jobqueue.Queue
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewJobManager(queue jobqueue.Queue) *jobManager {
	return &jobManager{
		queue: queue,
	}
}

func (m *jobManager) SubmitRequest(request *v1.JobRequest) (*v1.Job, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	request.SetDefaults()

	id := newID(request.Plotter)
	path, err := storeSVG(id, request)
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

func storeSVG(id string, request *v1.JobRequest) (string, error) {
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("could not create directory %s: %w", uploadDir, err)
	}

	path := path.Join(uploadDir, fmt.Sprintf("%s.svg", id))

	store, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("could not create file %s: %w", path, err)
	}
	defer store.Close()

	upload, err := request.SVG.Open()
	if err != nil {
		return "", fmt.Errorf("could not open file %s: %w", request.SVG.Filename, err)
	}
	defer upload.Close()

	size, err := io.Copy(store, upload)
	if err != nil {
		return "", fmt.Errorf("could not copy file %s: %w", request.SVG.Filename, err)
	}

	if size != request.SVG.Size {
		return "", errors.New("size mismatch")
	}

	return path, nil
}
