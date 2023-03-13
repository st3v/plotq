package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui/v4emb"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	v1 "github.com/st3v/plotq/api/v1"
)

const (
	tagJobs     = "Jobs"
	tagRequests = "JobRequests"
)

type Spooler interface {
	SubmitRequest(request *v1.JobRequest) (*v1.Job, error)
	GetJob(id string) (*v1.Job, error)
	GetJobs() ([]v1.Job, error)
	DeleteJob(id string) (*v1.Job, error)
}

func New(spooler Spooler) *web.Service {
	service := web.DefaultService()

	service.OpenAPI.Info.Title = "PlotterQueue API"
	service.OpenAPI.Info.WithDescription("Send job requests to HPGL plotters.")
	service.OpenAPI.Info.Version = "v1"

	service.Get("/v1/jobs", getJobs(spooler))
	service.Get("/v1/jobs/{id}", getJobByID(spooler))
	service.Post("/v1/jobs", postRequest(spooler))
	service.Delete("/v1/jobs/{id}", deleteJobByID(spooler))
	service.Docs("/v1/docs", v4emb.New)

	return service
}

func getJobs(spooler Spooler) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, _ struct{}, output *[]v1.Job) error {
		var err error
		*output, err = spooler.GetJobs()
		return err
	})

	u.SetTags(tagJobs)

	return u
}

func postRequest(spooler Spooler) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input v1.JobRequest, output *v1.Job) error {
		job, err := spooler.SubmitRequest(&input)
		if err != nil {
			return fmt.Errorf("failed to not submit request: %w", err)
		}

		*output = *job
		return nil
	})

	u.SetTags(tagRequests)
	u.SetExpectedErrors(status.AlreadyExists)

	return u
}

func getJobByID(spooler Spooler) usecase.Interactor {
	type idInput struct {
		ID string `path:"id" required:"true" example:"hp7550-5fbbd6p8"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input idInput, output *v1.Job) error {
		job, err := spooler.GetJob(input.ID)
		if err == nil && job == nil {
			return status.Wrap(errors.New("job not found"), status.NotFound)
		}

		*output = *job
		return err
	})

	u.SetTags(tagJobs)
	u.SetExpectedErrors(status.NotFound)

	return u
}

func deleteJobByID(spooler Spooler) usecase.Interactor {
	type idInput struct {
		ID string `path:"id" required:"true" example:"hp7550-5fbbd6p8"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input idInput, output *v1.Job) error {
		job, err := spooler.DeleteJob(input.ID)
		if err == nil && job == nil {
			return status.Wrap(errors.New("job not found"), status.NotFound)
		}

		*output = *job
		return err
	})

	u.SetTags(tagJobs)
	u.SetExpectedErrors(status.NotFound)

	return u
}
