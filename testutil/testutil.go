package testutil

import (
	"context"
	"math/rand"
	"net"
	"testing"
	"time"

	v1 "github.com/st3v/plotq/api/v1"
	"github.com/st3v/plotq/plotter"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandJob() v1.Job {
	return v1.Job{
		ID:      RandString(5),
		User:    RandString(5),
		Plotter: RandString(5),
		Settings: v1.JobSettings{
			Pagesize:    RandPagesize(),
			Velocity:    uint8(rand.Intn(100)),
			Orientation: RandOrientation(),
			Device:      RandDevice(),
		},
		SVG:         RandString(10),
		Status:      RandStatus(),
		SubmittedAt: time.Now(),
		Error:       RandString(10),
	}
}

const alphanumeric = "0123456789abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(b)
}

func RandPagesize() v1.Pagesize {
	s := v1.PagesizeA0.Enum()
	return s[rand.Intn(len(s))].(v1.Pagesize)
}

func RandOrientation() v1.Orientation {
	o := v1.OrientationPortrait.Enum()
	return o[rand.Intn(len(o))].(v1.Orientation)
}

func RandDevice() v1.Device {
	d := v1.DeviceArtisan.Enum()
	return d[rand.Intn(len(d))].(v1.Device)
}

func RandStatus() v1.JobStatus {
	s := v1.JobStatusPending.Enum()
	return s[rand.Intn(len(s))].(v1.JobStatus)
}

func NewTestServer(t *testing.T, expectedPayload []byte) *Testserver {
	server := &Testserver{
		T:               t,
		addr:            "localhost:3000",
		Ack:             "OK",
		Sleep:           0,
		ExpectedBufLen:  254,
		ExpectedPayload: expectedPayload,
	}

	server.ctx, server.cancel = context.WithCancel(context.Background())

	var err error
	server.listener, err = net.Listen("tcp", server.addr)
	require.NoError(t, err)

	return server
}

type Testserver struct {
	*testing.T
	ctx             context.Context
	cancel          context.CancelFunc
	listener        net.Listener
	addr            string
	Ack             string
	Sleep           time.Duration
	ExpectedPayload []byte
	ExpectedBufLen  int
}

func (t *Testserver) Addr() string {
	return t.addr
}

func (t *Testserver) Close() {
	t.cancel()
	t.listener.Close()
}

func (t *Testserver) MustConnect() *plotter.Conn {
	addr, err := t.acceptConnections()
	require.NoError(t, err)

	conn, err := plotter.Connect(addr, plotter.WithTimeout(time.Second))
	require.NoError(t, err)

	return conn
}

func (t *Testserver) acceptConnections() (string, error) {
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
					buf := make([]byte, t.ExpectedBufLen)

					n, err := conn.Read(buf)
					if err != nil {
						require.ErrorContains(t, err, "EOF")
						require.Equal(t, read, len(t.ExpectedPayload))
						return
					}

					// verify that the feeder does not write more than expected
					require.GreaterOrEqual(t, t.ExpectedBufLen, n)

					// verify that the feeder wrote the expected payload
					require.Equal(t, t.ExpectedPayload[read:read+n], buf[:n])
					read += n

					// sleep for a while to simulate a slow feeder
					time.Sleep(t.Sleep)

					// send ack to the feeder
					_, err = conn.Write([]byte(t.Ack))
					require.NoError(t, err)
				}
			}
		}
	}()

	return t.addr, nil
}
