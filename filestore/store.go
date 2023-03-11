package filestore

import "io"

type Store interface {
	Put(name string, src io.Reader) (written int64, err error)
	Get(name string) (file io.ReadCloser, err error)
}
