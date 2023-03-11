package converter

import "io"

type Converter interface {
	Convert(svg io.Reader, opts ...Option) (hpgl io.ReadCloser, err error)
}

type Option func(c *converterConfig)

func Landscape(c *converterConfig) {
	c.landscape = true
}

func Portrait(c *converterConfig) {
	c.landscape = false
}

func Pagesize(size string) Option {
	return func(c *converterConfig) {
		c.pagesize = size
	}
}

func Device(device string) Option {
	return func(c *converterConfig) {
		c.device = device
	}
}

func Velocity(velocity uint8) Option {
	return func(c *converterConfig) {
		c.velocity = velocity
	}
}

type converterConfig struct {
	pagesize  string
	device    string
	landscape bool
	velocity  uint8
}

func config(opts []Option) converterConfig {
	c := converterConfig{}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}
