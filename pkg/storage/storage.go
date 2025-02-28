package storage

import "io"

type MusicStorage interface {
	// GetMusic retorna um reader para a música
	GetMusic(id string) (io.ReadCloser, error)

	// SaveMusic salva uma música e retorna o caminho/id
	SaveMusic(id string, data io.Reader) error

	// DeleteMusic remove uma música
	DeleteMusic(id string) error
}
