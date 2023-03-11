package plotter

import (
	"time"
)

const DefaultTimeout = time.Minute

type configOptions struct {
	timeout time.Duration
}

type Option func(*configOptions)

func WithTimeout(d time.Duration) Option {
	return func(c *configOptions) {
		c.timeout = d
	}
}

func config(opts []Option) *configOptions {
	c := &configOptions{
		timeout: DefaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
