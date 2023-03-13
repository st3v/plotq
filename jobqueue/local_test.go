package jobqueue_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/beeker1121/goque"
	"github.com/stretchr/testify/require"

	v1 "github.com/st3v/plotq/api/v1"
	"github.com/st3v/plotq/jobqueue"
	"github.com/st3v/plotq/testutil"
)

func TestEnqueue(t *testing.T) {
	dir := t.TempDir()
	defer os.RemoveAll(dir)

	local, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)

	expected := testutil.RandJob()
	err = local.Enqueue(&expected)
	require.NoError(t, err)

	local.Close()

	queue, err := goque.OpenQueue(dir)
	require.NoError(t, err)
	defer queue.Close()

	item, err := queue.Dequeue()
	actual := &v1.Job{}
	err = json.Unmarshal(item.Value, actual)
	require.NoError(t, err)

	require.True(t, actual.SubmittedAt.Equal(expected.SubmittedAt))
	expected.SubmittedAt = actual.SubmittedAt
	require.Equal(t, expected, *actual)
}

func TestGetAll(t *testing.T) {
	dir := t.TempDir()
	defer os.RemoveAll(dir)

	local, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)
	defer local.Close()

	expected := make([]v1.Job, 10)
	for i := range expected {
		expected[i] = testutil.RandJob()
		err = local.Enqueue(&expected[i])
		require.NoError(t, err)
	}

	// we should be able to get the jobs multiple times
	for i := 0; i < 3; i++ {
		actual, err := local.GetAll()
		require.NoError(t, err)
		for i := range actual {
			require.True(t, actual[i].SubmittedAt.Equal(expected[i].SubmittedAt))
			expected[i].SubmittedAt = actual[i].SubmittedAt
			require.Equal(t, expected[i], actual[i])
		}
	}
}

func TestGet(t *testing.T) {
	dir := t.TempDir()
	defer os.RemoveAll(dir)

	local, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)
	defer local.Close()

	expected := make([]v1.Job, 10)
	for i := range expected {
		expected[i] = testutil.RandJob()
		err = local.Enqueue(&expected[i])
		require.NoError(t, err)
	}

	// we should be able to get the jobs multiple times
	for i := 0; i < 3; i++ {
		for i := len(expected) - 1; i >= 0; i-- {
			actual, err := local.Get(expected[i].ID)
			require.NoError(t, err)
			require.True(t, actual.SubmittedAt.Equal(expected[i].SubmittedAt))
			expected[i].SubmittedAt = actual.SubmittedAt
			require.Equal(t, expected[i], *actual)
		}
	}
}

func TestCancel(t *testing.T) {
	dir := t.TempDir()
	defer os.RemoveAll(dir)

	local, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)
	defer local.Close()

	expected := testutil.RandJob()
	expected.Status = v1.JobStatusPending
	err = local.Enqueue(&expected)
	require.NoError(t, err)

	actual, err := local.Get(expected.ID)
	require.NoError(t, err)
	require.Equal(t, v1.JobStatusPending, actual.Status)

	actual, err = local.Cancel(expected.ID)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, v1.JobStatusCanceled, actual.Status)

	actual, err = local.Get(expected.ID)
	require.NoError(t, err)
	require.Equal(t, v1.JobStatusCanceled, actual.Status)
}

func TestPeek(t *testing.T) {
	dir := t.TempDir()
	defer os.RemoveAll(dir)

	local, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)
	defer local.Close()

	expected := make([]v1.Job, 10)
	for i := range expected {
		expected[i] = testutil.RandJob()
		err = local.Enqueue(&expected[i])
		require.NoError(t, err)
	}

	for i := 0; i < 3; i++ {
		actual, err := local.Peek()
		require.NoError(t, err)
		require.True(t, actual.SubmittedAt.Equal(expected[0].SubmittedAt))
		expected[0].SubmittedAt = actual.SubmittedAt
		require.Equal(t, expected[0], *actual)
	}
}

func TestDequeue(t *testing.T) {
	dir := t.TempDir()
	defer os.RemoveAll(dir)

	local, err := jobqueue.OpenLocal(dir)
	require.NoError(t, err)
	defer local.Close()

	expected := make([]v1.Job, 10)
	for i := range expected {
		expected[i] = testutil.RandJob()
		err = local.Enqueue(&expected[i])
		require.NoError(t, err)
	}

	for i := range expected {
		actual, err := local.Dequeue()
		require.NoError(t, err)
		require.True(t, actual.SubmittedAt.Equal(expected[i].SubmittedAt))
		expected[i].SubmittedAt = actual.SubmittedAt
		require.Equal(t, expected[i], *actual)
	}

	actual, err := local.Dequeue()
	require.Error(t, jobqueue.ErrQueueEmpty)
	require.Nil(t, actual)
}
