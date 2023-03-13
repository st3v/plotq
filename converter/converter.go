package converter

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"io"

	v1 "github.com/st3v/plotq/api/v1"
)

// Convert returns an io.WriterTo that converts svg to hpgl
//
//counterfeiter:generate -o fake --fake-name Convert . Convert
type Convert func(svg io.Reader, opts ...Option) io.WriterTo

// Option is an option for a converter
type Option func(c *converterConfig)

// Orientation sets the orientation
func Orientation(orientation v1.Orientation) Option {
	return func(c *converterConfig) {
		if orientation == v1.OrientationLandscape {
			c.landscape = true
		}
	}
}

// Pagesize sets the pagesize
func Pagesize(size v1.Pagesize) Option {
	return func(c *converterConfig) {
		c.pagesize = string(size)
	}
}

// Device sets the device
func Device(device v1.Device) Option {
	return func(c *converterConfig) {
		c.device = string(device)
	}
}

// Velocity sets the velocity
func Velocity(velocity uint8) Option {
	return func(c *converterConfig) {
		c.velocity = velocity
	}
}

// converterConfig is the configuration for a converter
type converterConfig struct {
	pagesize  string
	device    string
	landscape bool
	velocity  uint8
}

// config returns a converterConfig from the given options
func config(opts []Option) converterConfig {
	c := converterConfig{}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}
