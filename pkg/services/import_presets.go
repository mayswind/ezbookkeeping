package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// ImportPresetsService represents import preset service
type ImportPresetsService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize an import preset service singleton instance
var (
	ImportPresets = &ImportPresetsService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetTotalImportPresetCountByUid returns total presets count of user
func (s *ImportPresetsService) GetTotalImportPresetCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.ImportPreset{})

	return count, err
}

// GetAllPresetNamesByUid returns all preset names
func (s *ImportPresetsService) GetAllPresetNamesByUid(c core.Context, uid int64) ([]*models.ImportPreset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var presets []*models.ImportPreset
	err := s.UserDataDB(uid).NewSession(c).Select("preset_id, uid, name, display_order").Where("uid=? AND deleted=?", uid, false).Find(&presets)

	return presets, err
}

// GetPresetByPresetId returns a preset model according to preset id
func (s *ImportPresetsService) GetPresetByPresetId(c core.Context, uid int64, presetId int64) (*models.ImportPreset, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if presetId <= 0 {
		return nil, errs.ErrImportPresetIdInvalid
	}

	preset := &models.ImportPreset{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(presetId).Where("uid=? AND deleted=?", uid, false).Get(preset)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrImportPresetNotFound
	}

	return preset, nil
}

// GetMaxDisplayOrder returns the max display order
func (s *ImportPresetsService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	preset := &models.ImportPreset{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(preset)

	if err != nil {
		return 0, err
	}

	if has {
		return preset.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreatePreset saves a new preset model to database
func (s *ImportPresetsService) CreatePreset(c core.Context, preset *models.ImportPreset) error {
	if preset.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	preset.PresetId = s.GenerateUuid(uuid.UUID_TYPE_IMPORT_PRESET)

	if preset.PresetId < 1 {
		return errs.ErrSystemIsBusy
	}

	preset.Deleted = false
	preset.CreatedUnixTime = time.Now().Unix()
	preset.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(preset.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(preset)
		return err
	})
}

// ModifyPreset saves an existed preset model to database
func (s *ImportPresetsService) ModifyPreset(c core.Context, preset *models.ImportPreset) error {
	if preset.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	preset.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(preset.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(preset.PresetId).Cols("name", "data", "updated_unix_time").Where("uid=? AND deleted=?", preset.Uid, false).Update(preset)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrImportPresetNotFound
		}

		return err
	})
}

// ModifyPresetDisplayOrders updates display order of given presets
func (s *ImportPresetsService) ModifyPresetDisplayOrders(c core.Context, uid int64, presets []*models.ImportPreset) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(presets); i++ {
		presets[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(presets); i++ {
			preset := presets[i]
			updatedRows, err := sess.ID(preset.PresetId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(preset)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrImportPresetNotFound
			}
		}

		return nil
	})
}

// DeletePreset deletes an existed preset from database
func (s *ImportPresetsService) DeletePreset(c core.Context, uid int64, presetId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.ImportPreset{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.ID(presetId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrImportPresetNotFound
		}

		return err
	})
}

// DeleteAllPresets deletes all existed presets from database
func (s *ImportPresetsService) DeleteAllPresets(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.ImportPreset{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		return nil
	})
}
