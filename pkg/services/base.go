package services

import (
	"fmt"
	"path/filepath"

	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/mail"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/storage"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// ServiceUsingDB represents a service that need to use db
type ServiceUsingDB struct {
	container *datastore.DataStoreContainer
}

// UserDB returns the datastore which contains user
func (s *ServiceUsingDB) UserDB() *datastore.Database {
	return s.container.UserStore.Choose(0)
}

// TokenDB returns the datastore which contains user token
func (s *ServiceUsingDB) TokenDB(uid int64) *datastore.Database {
	return s.container.TokenStore.Choose(uid)
}

// TokenDBByIndex returns the datastore which contains user token by index
func (s *ServiceUsingDB) TokenDBByIndex(index int) *datastore.Database {
	return s.container.TokenStore.Get(index)
}

// TokenDBCount returns the count of datastores which contains user token
func (s *ServiceUsingDB) TokenDBCount() int {
	return s.container.TokenStore.Count()
}

// UserDataDB returns the datastore which contains user data
func (s *ServiceUsingDB) UserDataDB(uid int64) *datastore.Database {
	return s.container.UserDataStore.Choose(uid)
}

// UserDataDBByIndex returns the datastore which contains user data by index
func (s *ServiceUsingDB) UserDataDBByIndex(index int) *datastore.Database {
	return s.container.UserDataStore.Get(index)
}

// UserDataDBCount returns the count of datastores which contains user data
func (s *ServiceUsingDB) UserDataDBCount() int {
	return s.container.UserDataStore.Count()
}

// ServiceUsingConfig represents a service that need to use config
type ServiceUsingConfig struct {
	container *settings.ConfigContainer
}

// CurrentConfig returns the current config
func (s *ServiceUsingConfig) CurrentConfig() *settings.Config {
	return s.container.Current
}

// ServiceUsingMailer represents a service that need to use mailer
type ServiceUsingMailer struct {
	container *mail.MailerContainer
}

// SendMail sends an email according to argument
func (s *ServiceUsingMailer) SendMail(message *mail.MailMessage) error {
	if s.container.Current == nil {
		return errs.ErrSMTPServerNotEnabled
	}

	return s.container.Current.SendMail(message)
}

// ServiceUsingUuid represents a service that need to use uuid
type ServiceUsingUuid struct {
	container *uuid.UuidContainer
}

// GenerateUuid generates a new uuid according to given uuid type
func (s *ServiceUsingUuid) GenerateUuid(uuidType uuid.UuidType) int64 {
	return s.container.GenerateUuid(uuidType)
}

// GenerateUuids generates new uuids according to given uuid type and count
func (s *ServiceUsingUuid) GenerateUuids(uuidType uuid.UuidType, count uint8) []int64 {
	return s.container.GenerateUuids(uuidType, count)
}

// ServiceUsingStorage represents a service that need to use storage
type ServiceUsingStorage struct {
	container *storage.StorageContainer
}

// ExistsAvatar returns whether the user avatar exists from the current avatar object storage
func (s *ServiceUsingStorage) ExistsAvatar(uid int64, fileExtension string) (bool, error) {
	return s.container.ExistsAvatar(s.getUserAvatarPath(uid, fileExtension))
}

// ReadAvatar returns the user avatar from the current avatar object storage
func (s *ServiceUsingStorage) ReadAvatar(uid int64, fileExtension string) (storage.ObjectInStorage, error) {
	return s.container.ReadAvatar(s.getUserAvatarPath(uid, fileExtension))
}

// SaveAvatar returns whether save the user avatar into the current avatar object storage successfully
func (s *ServiceUsingStorage) SaveAvatar(uid int64, object storage.ObjectInStorage, fileExtension string) error {
	return s.container.SaveAvatar(s.getUserAvatarPath(uid, fileExtension), object)
}

// DeleteAvatar returns whether delete the user avatar from the current avatar object storage successfully
func (s *ServiceUsingStorage) DeleteAvatar(uid int64, fileExtension string) error {
	return s.container.DeleteAvatar(s.getUserAvatarPath(uid, fileExtension))
}

// ExistsTransactionPicture returns whether the transaction picture exists from the current transaction picture object storage
func (s *ServiceUsingStorage) ExistsTransactionPicture(uid int64, pictureId int64, fileExtension string) (bool, error) {
	return s.container.ExistsTransactionPicture(s.getTransactionPicturePath(uid, pictureId, fileExtension))
}

// ReadTransactionPicture returns the transaction picture from the current transaction picture object storage
func (s *ServiceUsingStorage) ReadTransactionPicture(uid int64, pictureId int64, fileExtension string) (storage.ObjectInStorage, error) {
	return s.container.ReadTransactionPicture(s.getTransactionPicturePath(uid, pictureId, fileExtension))
}

// SaveTransactionPicture returns whether save the transaction picture into the current transaction picture object storage successfully
func (s *ServiceUsingStorage) SaveTransactionPicture(uid int64, pictureId int64, object storage.ObjectInStorage, fileExtension string) error {
	return s.container.SaveTransactionPicture(s.getTransactionPicturePath(uid, pictureId, fileExtension), object)
}

// DeleteTransactionPicture returns whether delete the transaction picture from the current transaction picture object storage successfully
func (s *ServiceUsingStorage) DeleteTransactionPicture(uid int64, pictureId int64, fileExtension string) error {
	return s.container.DeleteTransactionPicture(s.getTransactionPicturePath(uid, pictureId, fileExtension))
}

func (s *ServiceUsingStorage) getUserAvatarPath(uid int64, fileExtension string) string {
	return fmt.Sprintf("%d.%s", uid, fileExtension)
}

func (s *ServiceUsingStorage) getTransactionPicturePath(uid int64, pictureId int64, fileExtension string) string {
	return filepath.Join(utils.Int64ToString(uid), fmt.Sprintf("%d.%s", pictureId, fileExtension))
}
