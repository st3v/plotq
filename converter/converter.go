package converter

import "io"

// Converter converts svg to hpgl
type Converter interface {
	Convert(svg io.Reader, opts ...Option) io.WriterTo
}

// Option is an option for a converter
type Option func(c *converterConfig)

// Landscape sets the converter to landscape mode
func Landscape(c *converterConfig) {
	c.landscape = true
}

// Portrait sets the converter to portrait mode
func Portrait(c *converterConfig) {
	c.landscape = false
}

// Pagesize sets the pagesize
func Pagesize(size string) Option {
	return func(c *converterConfig) {
		c.pagesize = size
	}
}

// Device sets the device
func Device(device string) Option {
	return func(c *converterConfig) {
		c.device = device
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
