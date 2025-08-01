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

// GetTotalTransactionPicturesCountByUid returns total transaction pictures count of user
func (s *TransactionPictureService) GetTotalTransactionPicturesCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.TransactionPictureInfo{})

	return count, err
}

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

// GetNewPictureInfosByPictureIds returns new transaction picture info models according to transaction picture ids
func (s *TransactionPictureService) GetNewPictureInfosByPictureIds(c core.Context, uid int64, pictureIds []int64) ([]*models.TransactionPictureInfo, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if pictureIds == nil {
		return nil, errs.ErrTransactionPictureIdInvalid
	}

	var pictureInfos []*models.TransactionPictureInfo
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND transaction_id=?", uid, false, models.TransactionPictureNewPictureTransactionId).In("picture_id", pictureIds).OrderBy("picture_id asc").Find(&pictureInfos)

	if err != nil {
		return nil, err
	}

	return pictureInfos, nil
}

// GetPictureInfosByTransactionId returns transaction picture info models according to transaction id
func (s *TransactionPictureService) GetPictureInfosByTransactionId(c core.Context, uid int64, transactionId int64) ([]*models.TransactionPictureInfo, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if transactionId <= 0 {
		return nil, errs.ErrTransactionIdInvalid
	}

	var pictureInfos []*models.TransactionPictureInfo
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND transaction_id=?", uid, false, transactionId).OrderBy("picture_id asc").Find(&pictureInfos)

	if err != nil {
		return nil, err
	}

	return pictureInfos, nil
}

// GetPictureInfosByTransactionIds returns transaction picture info models according to transaction ids
func (s *TransactionPictureService) GetPictureInfosByTransactionIds(c core.Context, uid int64, transactionIds []int64) (map[int64][]*models.TransactionPictureInfo, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if transactionIds == nil {
		return nil, errs.ErrTransactionIdInvalid
	}

	var pictureInfos []*models.TransactionPictureInfo
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("transaction_id", transactionIds).OrderBy("picture_id asc").Find(&pictureInfos)

	if err != nil {
		return nil, err
	}

	pictureInfoMap := s.GetPictureInfoListMapByList(pictureInfos)
	return pictureInfoMap, err
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

	pictureFile, err := s.ReadTransactionPicture(c, pictureInfo.Uid, pictureInfo.PictureId, pictureInfo.PictureExtension)

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

	err := s.SaveTransactionPicture(c, pictureInfo.Uid, pictureInfo.PictureId, pictureFile, pictureInfo.PictureExtension)

	if err != nil {
		return err
	}

	return s.UserDataDB(pictureInfo.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(pictureInfo)
		return err
	})
}

// RemoveUnusedTransactionPicture removes the unused transaction picture of specified user
func (s *TransactionPictureService) RemoveUnusedTransactionPicture(c core.Context, uid int64, pictureId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if pictureId <= 0 {
		return errs.ErrTransactionPictureIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionPictureInfo{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDB().DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.ID(pictureId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=? AND transaction_id=?", uid, false, models.TransactionPictureNewPictureTransactionId).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionPictureNotFound
		}

		return err
	})
}

// GetPictureInfoMapByList returns a transaction picture info list map by a list
func (s *TransactionPictureService) GetPictureInfoMapByList(pictureInfos []*models.TransactionPictureInfo) map[int64]*models.TransactionPictureInfo {
	pictureInfoMap := make(map[int64]*models.TransactionPictureInfo)

	for i := 0; i < len(pictureInfos); i++ {
		pictureInfo := pictureInfos[i]
		pictureInfoMap[pictureInfo.PictureId] = pictureInfo
	}

	return pictureInfoMap
}

// GetPictureInfoListMapByList returns a transaction picture info list map by a list
func (s *TransactionPictureService) GetPictureInfoListMapByList(pictureInfos []*models.TransactionPictureInfo) map[int64][]*models.TransactionPictureInfo {
	pictureInfoMap := make(map[int64][]*models.TransactionPictureInfo)

	for i := 0; i < len(pictureInfos); i++ {
		pictureInfo := pictureInfos[i]

		pictureInfos, _ := pictureInfoMap[pictureInfo.TransactionId]
		pictureInfoMap[pictureInfo.TransactionId] = append(pictureInfos, pictureInfo)
	}

	return pictureInfoMap
}

// GetTransactionPictureIds returns transaction picture ids list
func (s *TransactionPictureService) GetTransactionPictureIds(pictureInfos []*models.TransactionPictureInfo) []int64 {
	pictureIds := make([]int64, len(pictureInfos))

	for i := 0; i < len(pictureInfos); i++ {
		pictureIds[i] = pictureInfos[i].PictureId
	}

	return pictureIds
}
