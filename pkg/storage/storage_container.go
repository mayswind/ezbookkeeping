package storage

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const avatarPathPrefix = "avatar"
const transactionPicturePathPrefix = "transaction"

// StorageContainer contains the current object storage
type StorageContainer struct {
	AvatarCurrentStorage             ObjectStorage
	TransactionPictureCurrentStorage ObjectStorage
}

// Initialize a object storage container singleton instance
var (
	Container = &StorageContainer{}
)

// InitializeStorageContainer initializes the current object storage according to the config
func InitializeStorageContainer(config *settings.Config) error {
	avatarStorage, err := newObjectStorage(config, avatarPathPrefix)

	if err != nil {
		return err
	}

	Container.AvatarCurrentStorage = avatarStorage

	transactionPictureStorage, err := newObjectStorage(config, transactionPicturePathPrefix)

	if err != nil {
		return err
	}

	Container.TransactionPictureCurrentStorage = transactionPictureStorage

	return nil
}

// ExistsAvatar returns whether the avatar file exists from the current avatar object storage
func (s *StorageContainer) ExistsAvatar(ctx core.Context, path string) (bool, error) {
	return s.AvatarCurrentStorage.Exists(ctx, path)
}

// ReadAvatar returns the avatar file from the current avatar object storage
func (s *StorageContainer) ReadAvatar(ctx core.Context, path string) (ObjectInStorage, error) {
	return s.AvatarCurrentStorage.Read(ctx, path)
}

// SaveAvatar returns whether save the avatar file into the current avatar object storage successfully
func (s *StorageContainer) SaveAvatar(ctx core.Context, path string, object ObjectInStorage) error {
	return s.AvatarCurrentStorage.Save(ctx, path, object)
}

// DeleteAvatar returns whether delete the avatar file from the current avatar object storage successfully
func (s *StorageContainer) DeleteAvatar(ctx core.Context, path string) error {
	return s.AvatarCurrentStorage.Delete(ctx, path)
}

// ExistsTransactionPicture returns whether the transaction picture file exists from the current transaction picture object storage
func (s *StorageContainer) ExistsTransactionPicture(ctx core.Context, path string) (bool, error) {
	return s.TransactionPictureCurrentStorage.Exists(ctx, path)
}

// ReadTransactionPicture returns the transaction picture file from the current transaction picture object storage
func (s *StorageContainer) ReadTransactionPicture(ctx core.Context, path string) (ObjectInStorage, error) {
	return s.TransactionPictureCurrentStorage.Read(ctx, path)
}

// SaveTransactionPicture returns whether save the transaction picture file into the current transaction picture object storage successfully
func (s *StorageContainer) SaveTransactionPicture(ctx core.Context, path string, object ObjectInStorage) error {
	return s.TransactionPictureCurrentStorage.Save(ctx, path, object)
}

// DeleteTransactionPicture returns whether delete the transaction picture file from the current transaction picture object storage successfully
func (s *StorageContainer) DeleteTransactionPicture(ctx core.Context, path string) error {
	return s.TransactionPictureCurrentStorage.Delete(ctx, path)
}

func newObjectStorage(config *settings.Config, pathPrefix string) (ObjectStorage, error) {
	if config.StorageType == settings.LocalFileSystemObjectStorageType {
		return NewLocalFileSystemObjectStorage(config, pathPrefix)
	} else if config.StorageType == settings.MinIOStorageType {
		return NewMinIOObjectStorage(config, pathPrefix)
	}

	return nil, errs.ErrInvalidStorageType
}
