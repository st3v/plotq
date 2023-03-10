package v1

import (
	"time"
)

type Job struct {
	ID          string      `json:"id" description:"ID is a unique string that identifies a job." example:"hp7550-5fbbd6p8"`
	User        string      `json:"user" description:"Name of the user that submitted the plot." example:"st3v"`
	Plotter     string      `json:"plotter" description:"Hostname of the plotter to use." example:"hp7550"`
	Settings    JobSettings `json:"settings" description:"Settings to use for the plot."`
	SVG         string      `json:"svg" description:"SVG file to be plotted." example:"uploads/hp7550-5fbbd6p8.svg"`
	HPGL        string      `json:"hpgl,omitempty" description:"Converted HPGL file." example:"plots/hp7550-5fbbd6p8.hpgl"`
	Status      JobStatus   `json:"status" description:"Current status of the job." example:"Pending"`
	SubmittedAt time.Time   `json:"submittedAt" description:"Time when the job was submitted."`
	Error       string      `json:"error,omitempty" description:"Error message if the job failed." example:""`
}

type JobSettings struct {
	Device      Device      `json:"device" description:"Device configuration." example:"hp7550"`
	Pagesize    Pagesize    `json:"pagesize" description:"Pagesize of plot." example:"a4"`
	Orientation Orientation `json:"orientation,ommitempty" description:"Orientation of plot." default:"portrait" example:"landscape"`
	Velocity    uint8       `json:"velocity,ommitempty" description:"Velocity to use for plotting." example:"50"`
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "Pending"
	JobStatusConverting JobStatus = "Converting"
	JobStatusPlotting   JobStatus = "Plotting"
	JobStatusCanceled   JobStatus = "Canceled"
	JobStatusSucceeded  JobStatus = "Succeeded"
	JobStatusFailed     JobStatus = "Failed"
)

func (JobStatus) Enum() []interface{} {
	return []interface{}{
		JobStatusPending,
		JobStatusConverting,
		JobStatusPlotting,
		JobStatusCanceled,
		JobStatusSucceeded,
		JobStatusFailed,
	}
}
