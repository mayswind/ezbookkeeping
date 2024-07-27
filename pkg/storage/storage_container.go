package storage

import (
	"fmt"

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
		storage, err := NewLocalFileSystemObjectStorage(config, avatarPathPrefix)
		Container.AvatarCurrentStorage = storage

		return err
	}

	return errs.ErrInvalidStorageType
}

// ExistsAvatar returns whether the user avatar exists from the current object storage
func (s *StorageContainer) ExistsAvatar(uid int64, fileExtension string) (bool, error) {
	return s.AvatarCurrentStorage.Exists(s.getUserAvatarPath(uid, fileExtension))
}

// ReadAvatar returns the user avatar from the current object storage
func (s *StorageContainer) ReadAvatar(uid int64, fileExtension string) (ObjectInStorage, error) {
	return s.AvatarCurrentStorage.Read(s.getUserAvatarPath(uid, fileExtension))
}

// SaveAvatar returns whether save the user avatar into the current object storage successfully
func (s *StorageContainer) SaveAvatar(uid int64, object ObjectInStorage, fileExtension string) error {
	return s.AvatarCurrentStorage.Save(s.getUserAvatarPath(uid, fileExtension), object)
}

// DeleteAvatar returns whether delete the user avatar from the current object storage successfully
func (s *StorageContainer) DeleteAvatar(uid int64, fileExtension string) error {
	return s.AvatarCurrentStorage.Delete(s.getUserAvatarPath(uid, fileExtension))
}

func (s *StorageContainer) getUserAvatarPath(uid int64, fileExtension string) string {
	return fmt.Sprintf("%d.%s", uid, fileExtension)
}
