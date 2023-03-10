package v1

import (
	"errors"
	"mime/multipart"
)

type JobRequest struct {
	User        string                `formData:"user" description:"Name of the user submitting the plot request." required:"true"`
	Plotter     string                `formData:"plotter" description:"Hostname of the plotter to use." required:"true" example:"hp7550"`
	Device      Device                `formData:"device" description:"Device configuration." required:"true"`
	Pagesize    Pagesize              `formData:"pagesize" description:"Pagesize of plot." required:"true"`
	Orientation Orientation           `formData:"orientation,omitempty" description:"Orientation of plot."`
	Velocity    uint8                 `formData:"velocity,omitempty" description:"Plotting velocity." example:"50"`
	SVG         *multipart.FileHeader `formData:"svg" description:"SVG file to be plotted." required:"true"`
}

func (r *JobRequest) Validate() error {
	if r.User == "" {
		return errors.New("no user specified")
	}

	if r.Plotter == "" {
		return errors.New("no plotter specified")
	}

	return nil
}

func (r *JobRequest) SetDefaults() {
	if r.Velocity == 0 {
		r.Velocity = DefaultVelocity
	}

	if r.Orientation == "" {
		r.Orientation = DefaultOrientation
	}
}
