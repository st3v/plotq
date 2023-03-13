package plotter_test

import (
	"context"
	"crypto/rand"
	"net"
	"testing"
	"time"

	"github.com/st3v/plotq/plotter"
	"github.com/stretchr/testify/require"
)

var hpgl = []byte("IN;DF;VS10;PS0;SP1;PA;PU0,10870;SP0;IN;\n")

func TestPlotterWrite(t *testing.T) {
	server := newTestServer(t, hpgl)
	defer server.close()

	conn := server.mustConnect()
	defer conn.Close()

	n, err := conn.Write(hpgl)
	require.NoError(t, err)
	require.Equal(t, len(hpgl), n)
}

func TestPlotterWriteLarge(t *testing.T) {
	twoMB := 2 * 1024 * 1024
	payload := make([]byte, twoMB)
	read, err := rand.Read(payload)
	require.NoError(t, err)
	require.Equal(t, twoMB, read)

	server := newTestServer(t, payload)
	defer server.close()

	conn := server.mustConnect()
	defer conn.Close()

	n, err := conn.Write(payload)
	require.NoError(t, err)
	require.Equal(t, len(payload), n)
}

func TestPlotterWriteEmpty(t *testing.T) {
	payload := make([]byte, 0)

	server := newTestServer(t, payload)
	defer server.close()

	conn := server.mustConnect()
	defer conn.Close()

	n, err := conn.Write(payload)
	require.NoError(t, err)
	require.Equal(t, len(payload), n)
}

func TestPlotterReadInvalidAck(t *testing.T) {
	server := newTestServer(t, hpgl)
	server.ack = "NO"
	defer server.close()

	conn := server.mustConnect()
	defer conn.Close()

	n, err := conn.Write(hpgl)
	require.ErrorContains(t, err, "did not ack with OK but NO")
	require.Equal(t, len(hpgl), n)
}

func TestPlotterTimeout(t *testing.T) {
	server := newTestServer(t, hpgl)
	server.sleep = 3 * time.Second
	defer server.close()

	conn := server.mustConnect()
	defer conn.Close()

	n, err := conn.Write(hpgl)
	require.ErrorContains(t, err, "timeout")
	require.Equal(t, len(hpgl), n)
}

func newTestServer(t *testing.T, expectedPayload []byte) *testserver {
	server := &testserver{
		T:               t,
		addr:            "localhost:3000",
		ack:             "OK",
		sleep:           0,
		expectedBufLen:  254,
		expectedPayload: expectedPayload,
	}

	server.ctx, server.cancel = context.WithCancel(context.Background())

	var err error
	server.listener, err = net.Listen("tcp", server.addr)
	require.NoError(t, err)

	return server
}

type testserver struct {
	*testing.T
	ctx             context.Context
	cancel          context.CancelFunc
	listener        net.Listener
	addr            string
	ack             string
	sleep           time.Duration
	expectedPayload []byte
	expectedBufLen  int
}

func (t *testserver) close() {
	t.cancel()
	t.listener.Close()
}

func (t *testserver) mustConnect() *plotter.Conn {
	addr, err := t.listen()
	require.NoError(t, err)

	conn, err := plotter.Connect(addr, plotter.WithTimeout(time.Second))
	require.NoError(t, err)

	return conn
}

func (t *testserver) listen() (string, error) {
	go func() {

		select {
		case <-t.ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
			conn, err := t.listener.Accept()
			require.NoError(t, err)
			defer conn.Close()

			read := 0
			for {
				select {
				case <-t.ctx.Done():
					return
				case <-time.After(time.Millisecond):
					buf := make([]byte, t.expectedBufLen)

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
			}
		}
	}()

	return t.addr, nil
}
