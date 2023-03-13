package spooler_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	v1 "github.com/st3v/plotq/api/v1"
	"github.com/st3v/plotq/filestore/fake"
	"github.com/st3v/plotq/jobqueue"
	"github.com/st3v/plotq/spooler"
	"github.com/st3v/plotq/testutil"
	"github.com/stretchr/testify/require"
)

func TestJobs(t *testing.T) {
	dir := t.TempDir()
	defer os.RemoveAll(dir)

	q, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)
	defer q.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s := spooler.NewSpooler(q, &fake.Store{})
	jobs := s.Incoming(ctx)

	expected := make([]v1.Job, 10)
	for i := range expected {
		expected[i] = testutil.RandJob()
		err = q.Enqueue(&expected[i])
		require.NoError(t, err)
	}

	var mu sync.Mutex
	actual := []v1.Job{}

	for loop := true; loop; {
		select {
		case <-ctx.Done():
			loop = false
		case job := <-jobs:
			mu.Lock()
			actual = append(actual, job)
			actualLen := len(actual)
			mu.Unlock()
			if actualLen == len(expected) {
				loop = false
			}
		}
	}

	require.Equal(t, len(expected), len(actual))

	for i := range actual {
		require.True(t, actual[i].SubmittedAt.Equal(expected[i].SubmittedAt))
		expected[i].SubmittedAt = actual[i].SubmittedAt
		require.Equal(t, expected[i], actual[i])
	}
}
