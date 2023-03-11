package filestore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type local struct {
	dir string
}

// local implements Store
var _ Store = &local{}

// NewLocalStore returns a new local file store
func NewLocalStore(dataDir string) (*local, error) {
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("could not create directory %s: %w", dataDir, err)
	}

	return &local{dir: dataDir}, nil
}

// Put writes the content of src to a file with the given name
func (l *local) Put(name string, src io.Reader) (int64, error) {
	path := filepath.Join(l.dir, name)

	dst, err := os.Create(path)
	if err != nil {
		return 0, fmt.Errorf("could not create file %s: %w", path, err)
	}

	written, err := io.Copy(dst, src)
	if err != nil {
		return 0, fmt.Errorf("could not copy from from src to file %w", err)
	}
	dst.Close()

	return written, nil
}

// Get returns a ReadCloser for the file with the given name
func (l *local) Get(name string) (io.ReadCloser, error) {
	path := filepath.Join(l.dir, name)
	return os.Open(path)
}
