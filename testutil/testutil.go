package testutil

import (
	"math/rand"
	"time"

	v1 "github.com/st3v/plotq/api/v1"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandJob() v1.Job {
	return v1.Job{
		ID:      RandString(5),
		User:    RandString(5),
		Plotter: RandString(5),
		Settings: v1.JobSettings{
			Pagesize:    RandPagesize(),
			Velocity:    uint8(rand.Intn(100)),
			Orientation: RandOrientation(),
			Device:      RandDevice(),
		},
		SVG:         RandString(10),
		Status:      RandStatus(),
		SubmittedAt: time.Now(),
		Error:       RandString(10),
	}
}

const alphanumeric = "0123456789abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

func RandPagesize() v1.Pagesize {
	s := v1.PagesizeA0.Enum()
	return s[rand.Intn(len(s))].(v1.Pagesize)
}

func RandOrientation() v1.Orientation {
	o := v1.OrientationPortrait.Enum()
	return o[rand.Intn(len(o))].(v1.Orientation)
}

func RandDevice() v1.Device {
	d := v1.DeviceArtisan.Enum()
	return d[rand.Intn(len(d))].(v1.Device)
}

func RandStatus() v1.JobStatus {
	s := v1.JobStatusPending.Enum()
	return s[rand.Intn(len(s))].(v1.JobStatus)
}
