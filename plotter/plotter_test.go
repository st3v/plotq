package plotter_test

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/st3v/plotq/testutil"
	"github.com/stretchr/testify/require"
)

var hpgl = []byte("IN;DF;VS10;PS0;SP1;PA;PU0,10870;SP0;IN;\n")

func TestPlotterWrite(t *testing.T) {
	server := testutil.NewTestServer(t, hpgl)
	defer server.Close()

	conn := server.MustConnect()
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

	server := testutil.NewTestServer(t, payload)
	defer server.Close()

	conn := server.MustConnect()
	defer conn.Close()

	n, err := conn.Write(payload)
	require.NoError(t, err)
	require.Equal(t, len(payload), n)
}

func TestPlotterWriteEmpty(t *testing.T) {
	payload := make([]byte, 0)

	server := testutil.NewTestServer(t, payload)
	defer server.Close()

	conn := server.MustConnect()
	defer conn.Close()

	n, err := conn.Write(payload)
	require.NoError(t, err)
	require.Equal(t, len(payload), n)
}

func TestPlotterReadInvalidAck(t *testing.T) {
	server := testutil.NewTestServer(t, hpgl)
	server.Ack = "NO"
	defer server.Close()

	conn := server.MustConnect()
	defer conn.Close()

	n, err := conn.Write(hpgl)
	require.ErrorContains(t, err, "did not ack with OK but NO")
	require.Equal(t, len(hpgl), n)
}

func TestPlotterTimeout(t *testing.T) {
	server := testutil.NewTestServer(t, hpgl)
	server.Sleep = 3 * time.Second
	defer server.Close()

	conn := server.MustConnect()
	defer conn.Close()

	n, err := conn.Write(hpgl)
	require.ErrorContains(t, err, "timeout")
	require.Equal(t, len(hpgl), n)
}
