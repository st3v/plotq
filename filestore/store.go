package filestore

import "io"

// Store is an interface for a file store
type Store interface {
	Put(name string, src io.Reader) (written int64, err error)
	Get(name string) (file io.ReadCloser, err error)
}
