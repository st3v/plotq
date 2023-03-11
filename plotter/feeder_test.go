package plotter_test

import (
	"bytes"
	"crypto/rand"
	"net"
	"testing"
	"time"

	"github.com/st3v/plotq/plotter"
	"github.com/stretchr/testify/require"
)

var hpgl = []byte("IN;DF;VS10;PS0;SP1;PA;PU0,10870;SP0;IN;\n")

func TestFeederPlot(t *testing.T) {
	server := newTestServer(t, hpgl)

	feeder := plotter.NewFeeder(server.Listen())
	err := feeder.Plot(bytes.NewReader(hpgl))
	require.NoError(t, err)
}

func TestFeederPlotLarge(t *testing.T) {
	fiveMB := 5 * 1024 * 1024
	payload := make([]byte, fiveMB)
	read, err := rand.Read(payload)
	require.NoError(t, err)
	require.Equal(t, fiveMB, read)

	server := newTestServer(t, payload)

	feeder := plotter.NewFeeder(server.Listen())
	err = feeder.Plot(bytes.NewReader(payload))
	require.NoError(t, err)
}

func TestFeederPlotEmpty(t *testing.T) {
	payload := make([]byte, 0)

	server := newTestServer(t, payload)

	feeder := plotter.NewFeeder(server.Listen())
	err := feeder.Plot(bytes.NewReader(payload))
	require.NoError(t, err)
}

func TestFeederPlotInvalidAck(t *testing.T) {
	server := newTestServer(t, hpgl)
	server.ack = "NOPE"

	feeder := plotter.NewFeeder(server.Listen())
	err := feeder.Plot(bytes.NewReader(hpgl))
	require.ErrorContains(t, err, "did not ack with OK but NOPE")
}

func TestFeederPlotTimeout(t *testing.T) {
	server := newTestServer(t, hpgl)
	server.sleep = 3 * time.Second

	feeder := plotter.NewFeeder(server.Listen(), plotter.WithTimeout(time.Second))
	err := feeder.Plot(bytes.NewReader(hpgl))
	require.ErrorContains(t, err, "timeout")
}

func newTestServer(t *testing.T, expectedPayload []byte) *testserver {
	return &testserver{
		T:               t,
		addr:            "localhost:3000",
		ack:             "OK",
		sleep:           0,
		expectedBufLen:  255,
		expectedPayload: expectedPayload,
	}
}

type testserver struct {
	*testing.T
	addr            string
	ack             string
	sleep           time.Duration
	expectedPayload []byte
	expectedBufLen  int
}

func (t *testserver) Listen() string {
	go func() {
		server, err := net.Listen("tcp", t.addr)
		require.NoError(t, err)
		defer server.Close()

		conn, err := server.Accept()
		require.NoError(t, err)
		defer conn.Close()

		read := 0
		for {
			buf := make([]byte, t.expectedBufLen*2)

			n, err := conn.Read(buf)
			if err != nil {
				require.ErrorContains(t, err, "EOF")
				require.Equal(t, read, len(t.expectedPayload))
				return
			}

			// verify that the feeder does not write more than expected
			require.GreaterOrEqual(t, t.expectedBufLen, n)

			// verify that the feeder wrote the expected payload
			require.Equal(t, t.expectedPayload[read:read+n], buf[:n])
			read += n

			// sleep for a while to simulate a slow feeder
			time.Sleep(t.sleep)

			// send ack to the feeder
			_, err = conn.Write([]byte(t.ack))
			require.NoError(t, err)
		}
	}()

	return t.addr
}
