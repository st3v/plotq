package plotter

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	bufLen = 255
	ack    = "OK"
)

type feeder struct {
	server  string
	timeout time.Duration
}

var _ io.Writer = &feeder{}

// NewFeeder creates a client for a PlotterFeeder server.
// See https://github.com/xHain-hackspace/PlotterFeeder
func NewFeeder(addr string, opts ...Option) *feeder {
	return &feeder{
		server:  addr,
		timeout: config(opts).timeout,
	}
}

// Plot sends the given HPGL data to the PlotterFeeder server.
func (f *feeder) Write(hpgl []byte) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", f.server)
	if err != nil {
		return 0, fmt.Errorf("could not resolve address %s: %w", f.server, err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return 0, fmt.Errorf("could not connect to %s: %w", addr, err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(f.timeout))

	reader := bytes.NewReader(hpgl)

	total := 0
	for err != io.EOF {
		n, err := io.CopyN(conn, reader, bufLen)
		if err != nil && err != io.EOF {
			return total, fmt.Errorf("could not copy to %s: %w", addr, err)
		} else if n == 0 {
			return total, nil
		}

		total += int(n)

		buf := make([]byte, bufLen)
		r, rerr := conn.Read(buf)
		if rerr != nil {
			return total, fmt.Errorf("could not read from %s: %w", addr, rerr)
		}

		got := string(buf[:r])
		if got != ack {
			return total, fmt.Errorf("server %s did not ack with %s but %s", addr, ack, got)
		}
	}

	return total, nil
}
