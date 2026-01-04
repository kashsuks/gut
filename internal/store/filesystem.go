package store

import (
	"os"
	"path/filepath"
)

type FileStore struct {
	RootPath string
}

func NewFileStore(root string) *FileStore {
	return &FileStore{RootPath: root}
}

func (s *FileStore) InitLayout() error {
	paths := []string{
		filepath.Join(s.RootPath, ".gut", "objects"),
		filepath.Join(s.RootPath, ".gut", "refs"),
	}

	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
		  return err
		}
	}

	return nil
}
