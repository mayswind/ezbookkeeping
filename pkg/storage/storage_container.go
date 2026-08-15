package storage

import (
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const avatarPathPrefix = "avatar"
const userCustomIconPathPrefix = "icon"
const transactionPicturePathPrefix = "transaction"

// StorageContainer contains the current object storage
type StorageContainer struct {
	avatarCurrentStorage             ObjectStorage
	userCustomIconCurrentStorage     ObjectStorage
	transactionPictureCurrentStorage ObjectStorage
}

// Initialize a object storage container singleton instance
var (
	Container = &StorageContainer{}
)

// InitializeStorageContainer initializes the current object storage according to the config
func InitializeStorageContainer(config *settings.Config) error {
	if config.AvatarProvider == core.USER_AVATAR_PROVIDER_INTERNAL {
		avatarStorage, err := newObjectStorage(config, avatarPathPrefix)

		if err != nil {
			return err
		}

		Container.avatarCurrentStorage = avatarStorage
	}

	if config.EnableUserCustomIcon {
		userCustomIconStorage, err := newObjectStorage(config, userCustomIconPathPrefix)

		if err != nil {
			return err
		}

		Container.userCustomIconCurrentStorage = userCustomIconStorage
	}

	if config.EnableTransactionPictures {
		transactionPictureStorage, err := newObjectStorage(config, transactionPicturePathPrefix)

		if err != nil {
			return err
		}

		Container.transactionPictureCurrentStorage = transactionPictureStorage
	}

	return nil
}

// ExistsAvatar returns whether the avatar file exists from the current avatar object storage
func (s *StorageContainer) ExistsAvatar(ctx core.Context, path string) (bool, error) {
	if s.avatarCurrentStorage == nil {
		return false, errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Exists(ctx, path)
}

// ReadAvatar returns the avatar file from the current avatar object storage
func (s *StorageContainer) ReadAvatar(ctx core.Context, path string) (ObjectInStorage, error) {
	if s.avatarCurrentStorage == nil {
		return nil, errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Read(ctx, path)
}

// SaveAvatar returns whether save the avatar file into the current avatar object storage successfully
func (s *StorageContainer) SaveAvatar(ctx core.Context, path string, object ObjectInStorage) error {
	if s.avatarCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Save(ctx, path, object)
}

// DeleteAvatar returns whether delete the avatar file from the current avatar object storage successfully
func (s *StorageContainer) DeleteAvatar(ctx core.Context, path string) error {
	if s.avatarCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.avatarCurrentStorage.Delete(ctx, path)
}

// ExistsUserCustomIcon returns whether the user custom icon file exists from the current user custom icon object storage
func (s *StorageContainer) ExistsUserCustomIcon(ctx core.Context, path string) (bool, error) {
	if s.userCustomIconCurrentStorage == nil {
		return false, errs.ErrSystemError
	}

	return s.userCustomIconCurrentStorage.Exists(ctx, path)
}

// ReadUserCustomIcon returns the user custom icon file from the current user custom icon object storage
func (s *StorageContainer) ReadUserCustomIcon(ctx core.Context, path string) (ObjectInStorage, error) {
	if s.userCustomIconCurrentStorage == nil {
		return nil, errs.ErrSystemError
	}

	return s.userCustomIconCurrentStorage.Read(ctx, path)
}

// SaveUserCustomIcon returns whether save the user custom icon file into the current user custom icon object storage successfully
func (s *StorageContainer) SaveUserCustomIcon(ctx core.Context, path string, object ObjectInStorage) error {
	if s.userCustomIconCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.userCustomIconCurrentStorage.Save(ctx, path, object)
}

// DeleteUserCustomIcon returns whether delete the user custom icon file from the current user custom icon object storage successfully
func (s *StorageContainer) DeleteUserCustomIcon(ctx core.Context, path string) error {
	if s.userCustomIconCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.userCustomIconCurrentStorage.Delete(ctx, path)
}

// ExistsTransactionPicture returns whether the transaction picture file exists from the current transaction picture object storage
func (s *StorageContainer) ExistsTransactionPicture(ctx core.Context, path string) (bool, error) {
	if s.transactionPictureCurrentStorage == nil {
		return false, errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Exists(ctx, path)
}

// ReadTransactionPicture returns the transaction picture file from the current transaction picture object storage
func (s *StorageContainer) ReadTransactionPicture(ctx core.Context, path string) (ObjectInStorage, error) {
	if s.transactionPictureCurrentStorage == nil {
		return nil, errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Read(ctx, path)
}

// SaveTransactionPicture returns whether save the transaction picture file into the current transaction picture object storage successfully
func (s *StorageContainer) SaveTransactionPicture(ctx core.Context, path string, object ObjectInStorage) error {
	if s.transactionPictureCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Save(ctx, path, object)
}

// DeleteTransactionPicture returns whether delete the transaction picture file from the current transaction picture object storage successfully
func (s *StorageContainer) DeleteTransactionPicture(ctx core.Context, path string) error {
	if s.transactionPictureCurrentStorage == nil {
		return errs.ErrSystemError
	}

	return s.transactionPictureCurrentStorage.Delete(ctx, path)
}

func newObjectStorage(config *settings.Config, pathPrefix string) (ObjectStorage, error) {
	if config.StorageType == settings.LocalFileSystemObjectStorageType {
		return NewLocalFileSystemObjectStorage(config, pathPrefix)
	} else if config.StorageType == settings.MinIOStorageType {
		return NewMinIOObjectStorage(config, pathPrefix)
	} else if config.StorageType == settings.WebDAVStorageType {
		return NewWebDAVObjectStorage(config, pathPrefix)
	}

	return nil, errs.ErrInvalidStorageType
}
