package spooler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	v1 "github.com/st3v/plotq/api/v1"
	"github.com/st3v/plotq/converter"
	"github.com/st3v/plotq/filestore"
	"github.com/st3v/plotq/jobqueue"
	"github.com/st3v/plotq/plotter"
)

const (
	// DefaultTick is the default tick duration for the spooler to check for the next job to be processed.
	DefaultTick = 1 * time.Second

	// DefaultTimeout is the default timeout for connections to the plotter.
	DefaultTimeout = time.Minute
)

type spooler struct {
	queue   jobqueue.Queue
	store   filestore.Store
	convert converter.Convert
	tick    time.Duration
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewSpooler creates a new job spooler.
func NewSpooler(queue jobqueue.Queue, svgStore filestore.Store, convert converter.Convert) *spooler {
	return &spooler{
		queue:   queue,
		store:   svgStore,
		convert: convert,
		tick:    DefaultTick,
	}
}

// SubmitRequest submits a new job request to the queue.
func (s *spooler) SubmitRequest(request *v1.JobRequest) (*v1.Job, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	request.SetDefaults()

	id := newID()
	path, err := s.storeSVG(id, request)
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

	if err := s.queue.Enqueue(job); err != nil {
		return nil, fmt.Errorf("failed to enqueue job: %w", err)
	}

	return job, nil
}

// GetJobs returns all jobs currently waiting to be processed.
func (s *spooler) GetJobs() ([]v1.Job, error) {
	return s.queue.GetAll()
}

// GetJob returns the job with the given ID from the queue of all jobs waiting to be processed.
func (s *spooler) GetJob(id string) (*v1.Job, error) {
	return s.queue.Get(id)
}

// DeleteJob deletes the job with the given ID.
func (s *spooler) DeleteJob(id string) (*v1.Job, error) {
	return s.queue.Cancel(id)
}

// Incoming returns a channel that receives jobs as they move to the front of the queue.
func (s *spooler) Incoming(ctx context.Context) <-chan v1.Job {
	jobs := make(chan v1.Job, 0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case <-time.After(s.tick):
				job, err := s.queue.Dequeue()
				if err == jobqueue.ErrQueueEmpty {
					continue
				} else if err != nil {
					log.Printf("failed to peek job: %v", err)
					continue
				}
				jobs <- *job
			}
		}
	}()
	return jobs
}

// Process processes a job.
func (s *spooler) Process(job v1.Job) (sent int64, err error) {
	file, err := s.store.Get(job.SVG)
	if err != nil {
		return 0, err
	}

	conn, err := plotter.Connect(job.Plotter, plotter.WithTimeout(DefaultTimeout))
	if err != nil {
		log.Printf("failed to connect to plotter %s: %v", job.Plotter, err)
		return 0, err
	}
	defer conn.Close()

	n, err := s.convert(
		file,
		converter.Orientation(job.Settings.Orientation),
		converter.Device(job.Settings.Device),
		converter.Velocity(job.Settings.Velocity),
		converter.Pagesize(job.Settings.Pagesize),
	).WriteTo(conn)

	if err != nil {
		return n, fmt.Errorf("failed to convert file and send to plotter: %v", err)
	}

	return n, nil
}

func (s *spooler) storeSVG(id string, request *v1.JobRequest) (string, error) {
	svg, err := request.SVG.Open()
	if err != nil {
		return "", fmt.Errorf("could not open file %s: %w", request.SVG.Filename, err)
	}
	defer svg.Close()

	path := fmt.Sprintf("%s.svg", id)

	written, err := s.store.Put(path, svg)
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

func newID() string {
	return randString(16)
}
