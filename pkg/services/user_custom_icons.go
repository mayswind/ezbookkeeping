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
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// UserCustomIconService represents user custom icon service
type UserCustomIconService struct {
	ServiceUsingDB
	ServiceUsingUuid
	ServiceUsingStorage
}

// Initialize a user custom icon service singleton instance
var (
	UserCustomIcons = &UserCustomIconService{
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

// GetTotalCustomIconsCountByUid returns total custom icons count of user
func (s *UserCustomIconService) GetTotalCustomIconsCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.UserCustomIcon{})

	return count, err
}

// GetAllCustomIconInfosByUid returns all custom icons of specified user
func (s *UserCustomIconService) GetAllCustomIconInfosByUid(c core.Context, uid int64) ([]*models.UserCustomIcon, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var customIcons []*models.UserCustomIcon
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Find(&customIcons)

	return customIcons, err
}

// GetCustomIconInfoByIconId returns a custom icon according to custom icon id
func (s *UserCustomIconService) GetCustomIconInfoByIconId(c core.Context, uid int64, iconId int64) (*models.UserCustomIcon, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if iconId <= 0 {
		return nil, errs.ErrUserCustomIconIdInvalid
	}

	customIconInfo := &models.UserCustomIcon{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(iconId).Where("uid=? AND deleted=?", uid, false).Get(customIconInfo)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserCustomIconNotFound
	}

	return customIconInfo, nil
}

// GetCustomIcon returns custom icon data according to custom icon id
func (s *UserCustomIconService) GetCustomIconByIconId(c core.Context, uid int64, iconId int64) ([]byte, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if iconId <= 0 {
		return nil, errs.ErrUserCustomIconIdInvalid
	}

	customIconInfo := &models.UserCustomIcon{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(iconId).Where("uid=? AND deleted=?", uid, false).Get(customIconInfo)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrUserCustomIconNotFound
	}

	customIconFile, err := s.ReadUserCustomIcon(c, uid, iconId)

	if os.IsNotExist(err) {
		return nil, errs.ErrUserCustomIconeNotExists
	}

	if err != nil {
		return nil, err
	}

	defer customIconFile.Close()

	customIconData, err := io.ReadAll(customIconFile)

	if err != nil {
		return nil, err
	}

	return customIconData, nil
}

// GetMaxDisplayOrder returns the max display order
func (s *UserCustomIconService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	customIcon := &models.UserCustomIcon{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(customIcon)

	if err != nil {
		return 0, err
	}

	if has {
		return customIcon.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// UploadCustomIcon uploads a custom icon for specified user
func (s *UserCustomIconService) UploadCustomIcon(c core.Context, customIcon *models.UserCustomIcon, customIconFile multipart.File) error {
	if customIcon.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	defer customIconFile.Close()

	customIcon.IconId = s.GenerateUuid(uuid.UUID_TYPE_CUSTOM_ICON)

	if customIcon.IconId < 1 {
		return errs.ErrSystemIsBusy
	}

	customIcon.Deleted = false
	customIcon.CreatedUnixTime = time.Now().Unix()
	customIcon.UpdatedUnixTime = time.Now().Unix()

	err := s.SaveUserCustomIcon(c, customIcon.Uid, customIcon.IconId, customIconFile)

	if err != nil {
		return err
	}

	return s.UserDataDB(customIcon.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(customIcon)
		return err
	})
}

// ModifyCustomIconDisplayOrders modifies display order of specified custom icons
func (s *UserCustomIconService) ModifyCustomIconDisplayOrders(c core.Context, uid int64, customIcons []*models.UserCustomIcon) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(customIcons); i++ {
		customIcons[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(customIcons); i++ {
			customIcon := customIcons[i]
			updatedRows, err := sess.ID(customIcon.IconId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(customIcon)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrUserCustomIconNotFound
			}
		}

		return nil
	})
}

// DeleteCustomIcon deletes specified custom icon
func (s *UserCustomIconService) DeleteCustomIcon(c core.Context, uid int64, iconId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.UserCustomIcon{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		accountUsed, err := sess.Where("uid=? AND deleted=? AND icon_type=? AND icon=?", uid, false, core.ICON_TYPE_USER_CUSTOM, iconId).Exist(&models.Account{})

		if err != nil {
			return err
		}

		if accountUsed {
			return errs.ErrUserCustomIconInUse
		}

		categoryUsed, err := sess.Where("uid=? AND deleted=? AND icon_type=? AND icon=?", uid, false, core.ICON_TYPE_USER_CUSTOM, iconId).Exist(&models.TransactionCategory{})

		if err != nil {
			return err
		}

		if categoryUsed {
			return errs.ErrUserCustomIconInUse
		}

		deletedRows, err := sess.ID(iconId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrUserCustomIconNotFound
		}

		return nil
	})
}

// ExistsCustomIcon returns whether the given user has existed custom icon
func (s *UserCustomIconService) ExistsCustomIcon(c core.Context, uid int64, iconId int64) (bool, error) {
	if uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	if iconId <= 0 {
		return false, errs.ErrUserCustomIconIdInvalid
	}

	return s.UserDataDB(uid).NewSession(c).ID(iconId).Where("uid=? AND deleted=?", uid, false).Exist(&models.UserCustomIcon{})
}

// ExistsCustomIcons returns whether the given user has existed custom icons
func (s *UserCustomIconService) ExistsCustomIcons(c core.Context, uid int64, iconIds []int64) (bool, error) {
	if uid <= 0 {
		return false, errs.ErrUserIdInvalid
	}

	if len(iconIds) == 0 {
		return false, errs.ErrUserCustomIconIdInvalid
	}

	uniqueIconIds := utils.ToUniqueInt64Slice(iconIds)
	count, err := s.UserDataDB(uid).NewSession(c).In("icon_id", uniqueIconIds).Where("uid=? AND deleted=?", uid, false).Count(&models.UserCustomIcon{})

	if err != nil {
		return false, err
	}

	return count == int64(len(uniqueIconIds)), nil
}
