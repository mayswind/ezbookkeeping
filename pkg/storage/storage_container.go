package storage

import (
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const avatarPathPrefix = "avatar"

// StorageContainer contains the current object storage
type StorageContainer struct {
	AvatarCurrentStorage ObjectStorage
}

// Initialize a object storage container singleton instance
var (
	Container = &StorageContainer{}
)

// InitializeStorageContainer initializes the current object storage according to the config
func InitializeStorageContainer(config *settings.Config) error {
	if config.StorageType == settings.LocalFileSystemObjectStorageType {
		avatarStorage, err := NewLocalFileSystemObjectStorage(config, avatarPathPrefix)
		Container.AvatarCurrentStorage = avatarStorage

		return err
	} else if config.StorageType == settings.MinIOStorageType {
		avatarStorage, err := NewMinIOObjectStorage(config, avatarPathPrefix)
		Container.AvatarCurrentStorage = avatarStorage

		return err
	}

	return errs.ErrInvalidStorageType
}

// ExistsAvatar returns whether the avatar file exists from the current avatar object storage
func (s *StorageContainer) ExistsAvatar(path string) (bool, error) {
	return s.AvatarCurrentStorage.Exists(path)
}

// ReadAvatar returns the avatar file from the current avatar object storage
func (s *StorageContainer) ReadAvatar(path string) (ObjectInStorage, error) {
	return s.AvatarCurrentStorage.Read(path)
}

// SaveAvatar returns whether save the avatar file into the current avatar object storage successfully
func (s *StorageContainer) SaveAvatar(path string, object ObjectInStorage) error {
	return s.AvatarCurrentStorage.Save(path, object)
}

// DeleteAvatar returns whether delete the avatar file from the current avatar object storage successfully
func (s *StorageContainer) DeleteAvatar(path string) error {
	return s.AvatarCurrentStorage.Delete(path)
}
