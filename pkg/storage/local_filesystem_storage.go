package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// LocalFileSystemObjectStorage represents local file system object storage
type LocalFileSystemObjectStorage struct {
	rootPath string
}

// NewLocalFileSystemObjectStorage returns a local file system object storage
func NewLocalFileSystemObjectStorage(config *settings.Config, pathPrefix string) (*LocalFileSystemObjectStorage, error) {
	storage := &LocalFileSystemObjectStorage{
		rootPath: filepath.Join(config.LocalFileSystemPath, pathPrefix),
	}

	if err := os.MkdirAll(storage.rootPath, os.ModePerm); err != nil {
		return nil, err
	}

	return storage, nil
}

// Exists returns whether the file exists
func (s *LocalFileSystemObjectStorage) Exists(ctx core.Context, path string) (bool, error) {
	return utils.IsExists(s.getFinalPath(path))
}

// Read returns the object instance according to specified the file path
func (s *LocalFileSystemObjectStorage) Read(ctx core.Context, path string) (ObjectInStorage, error) {
	return os.Open(s.getFinalPath(path))
}

// Save returns whether save the object instance successfully
func (s *LocalFileSystemObjectStorage) Save(ctx core.Context, path string, object ObjectInStorage) error {
	finalPath := s.getFinalPath(path)

	if err := os.MkdirAll(filepath.Dir(finalPath), os.ModePerm); err != nil {
		return err
	}

	targetFile, err := os.Create(finalPath)

	if err != nil {
		return err
	}

	defer targetFile.Close()

	_, err = io.Copy(targetFile, object)

	return err
}

// Delete returns whether delete the object according to specified the file path successfully
func (s *LocalFileSystemObjectStorage) Delete(ctx core.Context, path string) error {
	return os.Remove(s.getFinalPath(path))
}

func (s *LocalFileSystemObjectStorage) getFinalPath(path string) string {
	return filepath.Join(s.rootPath, path)
}
