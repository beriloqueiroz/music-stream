package storage

import "io"

type Storage interface {
	GetItem(id string) (io.ReadCloser, error)

	SaveItem(id string, data io.Reader) error

	DeleteItem(id string) error
}
