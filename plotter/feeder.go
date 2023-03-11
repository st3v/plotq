package plotter

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	writeBufLen = 255
	ack         = "OK"
)

type feeder struct {
	server  string
	timeout time.Duration
}

// NewFeeder creates a client for a PlotterFeeder server.
// See https://github.com/xHain-hackspace/PlotterFeeder
func NewFeeder(addr string, opts ...Option) *feeder {
	return &feeder{
		server:  addr,
		timeout: config(opts).timeout,
	}
}

// Plot sends the given HPGL data to the PlotterFeeder server.
func (f *feeder) Plot(hpgl io.Reader) error {
	addr, err := net.ResolveTCPAddr("tcp", f.server)
	if err != nil {
		return fmt.Errorf("could not resolve address %s: %w", f.server, err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return fmt.Errorf("could not connect to %s: %w", addr, err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(f.timeout))

	for {
		buf := make([]byte, writeBufLen)

		r, err := hpgl.Read(buf)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("could not read from hpgl: %w", err)
		}

		w, err := conn.Write(buf[:r])
		if err != nil {
			return fmt.Errorf("could not write to %s: %w", addr, err)
		}

		if w != r {
			return fmt.Errorf("could not write all bytes to %s: %d of %d", addr, w, r)
		}

		r, err = conn.Read(buf)
		if err != nil {
			return fmt.Errorf("could not read from %s: %w", addr, err)
		}

		got := string(buf[:r])
		if got != ack {
			return fmt.Errorf("server %s did not ack with %s but %s", addr, ack, got)
		}
	}
}
