package store

import (
	"os"
	"compress/zlib"
	"gut/pkg/models"
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

	headPath := filepath.Join(s.RootPath, ".gut", "HEAD")
	return os.WriteFile(headPath, []byte("ref: refs/main\n"), 0644)
}

func (s *FileStore) WriteObject(obj models.GutObject) error {
	dir := filepath.Join(s.RootPath, ".gut", "objects", obj.HashSum[:2])
	path := filepath.Join(dir, obj.HashSum[2:])

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	zw := zlib.NewWriter(f)
	defer zw.Close()

	_, err = zw.Write(obj.Content)
	return err
}
