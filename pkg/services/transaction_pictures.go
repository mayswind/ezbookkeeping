package services

import (
	"io"
	"mime/multipart"
	"os"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/storage"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// TransactionPictureService represents transaction picture service
type TransactionPictureService struct {
	ServiceUsingDB
	ServiceUsingUuid
	ServiceUsingStorage
}

// Initialize a transaction picture service singleton instance
var (
	TransactionPictures = &TransactionPictureService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
		ServiceUsingStorage: ServiceUsingStorage{
			container: storage.Container,
		},
	}
)

// GetPictureInfoByPictureId returns a transaction picture info model according to transaction picture id
func (s *TransactionPictureService) GetPictureInfoByPictureId(c core.Context, uid int64, pictureId int64) (*models.TransactionPictureInfo, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if pictureId <= 0 {
		return nil, errs.ErrTransactionPictureIdInvalid
	}

	pictureInfo := &models.TransactionPictureInfo{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(pictureId).Where("uid=? AND deleted=?", uid, false).Get(pictureInfo)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionPictureNotFound
	}

	return pictureInfo, nil
}

// GetPictureByPictureId returns the transaction picture data according to transaction picture id
func (s *TransactionPictureService) GetPictureByPictureId(c core.Context, uid int64, pictureId int64, fileExtension string) ([]byte, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if pictureId <= 0 {
		return nil, errs.ErrTransactionPictureIdInvalid
	}

	pictureInfo := &models.TransactionPictureInfo{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(pictureId).Where("uid=? AND deleted=?", uid, false).Get(pictureInfo)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionPictureNotFound
	}

	if pictureInfo.PictureExtension == "" {
		return nil, errs.ErrTransactionPictureNotFound
	}

	if pictureInfo.PictureExtension != fileExtension {
		return nil, errs.ErrTransactionPictureExtensionInvalid
	}

	pictureFile, err := s.ReadTransactionPicture(pictureInfo.Uid, pictureInfo.PictureId, pictureInfo.PictureExtension)

	if os.IsNotExist(err) {
		return nil, errs.ErrTransactionPictureNoExists
	}

	if err != nil {
		return nil, err
	}

	defer pictureFile.Close()

	pictureData, err := io.ReadAll(pictureFile)

	if err != nil {
		return nil, err
	}

	return pictureData, nil
}

// UploadPicture uploads the transaction picture for specified user
func (s *TransactionPictureService) UploadPicture(c core.Context, pictureInfo *models.TransactionPictureInfo, pictureFile multipart.File) error {
	if pictureInfo.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	defer pictureFile.Close()

	pictureInfo.PictureId = s.GenerateUuid(uuid.UUID_TYPE_USER)

	if pictureInfo.PictureId < 1 {
		return errs.ErrSystemIsBusy
	}

	pictureInfo.TransactionId = 0
	pictureInfo.Deleted = false
	pictureInfo.CreatedUnixTime = time.Now().Unix()
	pictureInfo.UpdatedUnixTime = time.Now().Unix()

	err := s.SaveTransactionPicture(pictureInfo.Uid, pictureInfo.PictureId, pictureFile, pictureInfo.PictureExtension)

	if err != nil {
		return err
	}

	return s.UserDataDB(pictureInfo.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(pictureInfo)
		return err
	})
}
