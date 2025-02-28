package local

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

func (s *LocalStorage) SaveMusic(id string, data io.Reader) error {
	path := filepath.Join(s.basePath, id)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

func (s *LocalStorage) GetMusic(id string) (io.ReadCloser, error) {
	path := filepath.Join(s.basePath, id)
	return os.Open(path)
}

func (s *LocalStorage) DeleteMusic(id string) error {
	path := filepath.Join(s.basePath, id)
	return os.Remove(path)
}
