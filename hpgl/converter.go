package hpgl

import "io"

type Converter interface {
	Convert(svg io.Reader, opts ...ConvertOption) (hpgl io.ReadCloser, err error)
}

type ConvertOption func(c *converterConfig)

func Landscape(c *converterConfig) {
	c.landscape = true
}

func Portrait(c *converterConfig) {
	c.landscape = false
}

func Pagesize(size string) ConvertOption {
	return func(c *converterConfig) {
		c.pagesize = size
	}
}

func Device(device string) ConvertOption {
	return func(c *converterConfig) {
		c.device = device
	}
}

func Velocity(velocity uint8) ConvertOption {
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

func config(opts []ConvertOption) converterConfig {
	c := converterConfig{}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}
