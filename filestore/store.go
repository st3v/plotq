package filestore

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import "io"

// Store is an interface for a file store
//
//counterfeiter:generate -o fake --fake-name Store . Store
type Store interface {
	Put(name string, src io.Reader) (written int64, err error)
	Get(name string) (file io.ReadCloser, err error)
}
