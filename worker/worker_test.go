package worker_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	v1 "github.com/st3v/plotq/api/v1"
	"github.com/st3v/plotq/converter"
	converterfake "github.com/st3v/plotq/converter/fake"
	"github.com/st3v/plotq/filestore"
	filestorefake "github.com/st3v/plotq/filestore/fake"
	"github.com/st3v/plotq/jobqueue"
	"github.com/st3v/plotq/plotter"
	"github.com/st3v/plotq/spooler"
	"github.com/st3v/plotq/testutil"
	"github.com/st3v/plotq/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	t.Skip("skipping e2e test")
	store, err := filestore.NewLocalStore(filepath.Join("..", "data", "upload"))
	require.NoError(t, err)

	file, err := store.Get("hp7550-vdjy34dq.svg")
	require.NoError(t, err)

	conn, err := plotter.Connect("hp-7550:1337", plotter.WithTimeout(time.Minute))
	require.NoError(t, err)
	defer conn.Close()

	vpype := converter.Vpype()
	n, err := vpype.Convert(file,
		converter.Orientation(v1.OrientationLandscape),
		converter.Device("hp7550"),
		converter.Velocity(10),
		converter.Pagesize("a3"),
	).WriteTo(conn)
	assert.Equal(t, int64(1318242), n)
	require.NoError(t, err)
}

func TestWorker(t *testing.T) {
	svg := `<svg height="50" width="50"><line x1="0" y1="0" x2="50" y2="50" style="stroke:black"/></svg>`
	hpgl := []byte("IN;DF;VS10;PS0;SP1;PA;PU0,10870;SP0;IN;\n")
	expected := []byte("huh")

	files := &filestorefake.Store{}
	files.GetReturns(io.NopCloser(strings.NewReader(svg)), nil)

	convert := &converterfake.Convert{}
	buf := bytes.NewBuffer(hpgl)
	convert.Returns(buf)

	timeout := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	plotter := testutil.NewTestServer(t, expected)
	defer plotter.Close()

	job := testutil.RandJob()
	job.Plotter = plotter.Addr()

	dir := t.TempDir()
	defer os.RemoveAll(dir)

	queue, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)
	defer queue.Close()

	err = queue.Enqueue(&job)
	require.NoError(t, err)

	spooler := spooler.NewSpooler(queue, files, convert.Spy)

	go worker.Run(ctx, spooler)

	<-time.After(timeout)

	require.Equal(t, 1, files.GetCallCount())
	require.Equal(t, 1, convert.CallCount())
}
