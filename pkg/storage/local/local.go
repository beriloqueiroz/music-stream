package local

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{
		baseDir: baseDir,
	}
}

func (s *LocalStorage) GetMusic(id string) (io.ReadCloser, error) {
	path := filepath.Join(s.baseDir, id+".mp3")
	return os.Open(path)
}

func (s *LocalStorage) SaveMusic(id string, data io.Reader) error {
	path := filepath.Join(s.baseDir, id+".mp3")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

func (s *LocalStorage) DeleteMusic(id string) error {
	path := filepath.Join(s.baseDir, id+".mp3")
	return os.Remove(path)
}
