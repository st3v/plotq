package plotter

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	// defaultTimeout is the default timeout for reads and writesto the PlotterFeeder
	defaultTimeout = time.Minute

	// buflen is the buffer size for writes to the PlotterFeeder
	buflen = 254

	// ack is the expected ack message from the PlotterFeeder
	ack = "OK"
)

// Option is a configuration option .
type ConnOption func(*connOptions)

// WithTimeout sets the timeout for the connection.
func WithTimeout(d time.Duration) ConnOption {
	return func(c *connOptions) {
		c.timeout = d
	}
}

// Conn represents a connection to a PlotterFeeder.
type Conn struct {
	conn    net.Conn
	timeout time.Duration
}

// feed implements io.WriteCloser.
var _ io.WriteCloser = &Conn{}

// connOptions is the configuration for a connection.
type connOptions struct {
	timeout time.Duration
}

// Connect creates a new connection to a PlotterFeeder.
// See https://github.com/xHain-hackspace/PlotterFeeder
func Connect(addr string, opts ...ConnOption) (*Conn, error) {
	server, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve address %s: %w", server, err)
	}

	conn, err := net.DialTCP("tcp", nil, server)
	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %w", addr, err)
	}

	return &Conn{
		conn:    conn,
		timeout: config(opts).timeout,
	}, nil
}

// Close closes the connection to the PlotterFeeder.
func (c *Conn) Close() error {
	fmt.Println("closing conn")
	return c.conn.Close()
}

// Plot sends the given HPGL data to the PlotterFeeder server.
func (c *Conn) Write(hpgl []byte) (int, error) {
	var err error
	reader := bytes.NewReader(hpgl)
	total := 0
	for err != io.EOF {
		c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
		n, err := io.CopyN(c.conn, reader, buflen)
		if err != nil && err != io.EOF {
			return total, fmt.Errorf("could not copy to server: %w", err)
		} else if n == 0 {
			return total, nil
		}

		total += int(n)

		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		buf := make([]byte, len(ack))
		_, rerr := c.conn.Read(buf)
		if rerr != nil {
			return total, fmt.Errorf("could not read from sever: %w", rerr)
		}

		got := string(buf)
		if got != ack {
			return total, fmt.Errorf("server did not ack with %s but %s", ack, got)
		}
	}

	return total, nil
}

// config creates new connOptions
func config(opts []ConnOption) *connOptions {
	c := &connOptions{
		timeout: defaultTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
