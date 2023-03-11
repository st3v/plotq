package filestore_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/st3v/plotq/filestore"
	"github.com/stretchr/testify/require"
)

func TestLocalPut(t *testing.T) {
	dir := t.TempDir()
	name := "foo"
	contents := "bar"

	local, err := filestore.NewLocalStore(dir)
	require.NoError(t, err)

	n, err := local.Put(name, strings.NewReader(contents))
	require.NoError(t, err)
	require.Equal(t, int64(len(contents)), n)

	expectedFile := filepath.Join(dir, name)
	require.FileExists(t, expectedFile)
	f, err := os.Open(expectedFile)
	require.NoError(t, err)
	defer f.Close()

	actualContents, err := io.ReadAll(f)
	require.NoError(t, err)
	require.Equal(t, contents, string(actualContents))
}

func TestLocalGet(t *testing.T) {
	dir := t.TempDir()
	name := "foo"
	contents := "bar"

	local, err := filestore.NewLocalStore(dir)
	require.NoError(t, err)

	f, err := os.Create(filepath.Join(dir, name))
	require.NoError(t, err)
	defer f.Close()
	defer os.Remove(f.Name())

	n, err := f.WriteString(contents)
	require.NoError(t, err)
	require.Equal(t, len(contents), n)

	reader, err := local.Get(name)
	require.NoError(t, err)
	defer reader.Close()

	actualContents, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, contents, string(actualContents))
}

func TestLocalRoundtrip(t *testing.T) {
	dir := t.TempDir()
	name := "foo"
	contents := "bar"

	local, err := filestore.NewLocalStore(dir)
	require.NoError(t, err)

	w, err := local.Put(name, strings.NewReader(contents))
	require.NoError(t, err)
	require.Equal(t, int64(len(contents)), w)

	reader, err := local.Get(name)
	require.NoError(t, err)
	defer reader.Close()

	actualContents, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, contents, string(actualContents))
}
