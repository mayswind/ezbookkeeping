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

// InsightsExplorerService represents insights explorer service
type InsightsExplorerService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a insights explorer service singleton instance
var (
	InsightsExplorers = &InsightsExplorerService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetTotalInsightsExplorersCountByUid returns total insights explorers count of user
func (s *InsightsExplorerService) GetTotalInsightsExplorersCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.InsightsExplorer{})

	return count, err
}

// GetAllInsightsExplorerNamesByUid returns all insights explorer models of user without data
func (s *InsightsExplorerService) GetAllInsightsExplorerNamesByUid(c core.Context, uid int64) ([]*models.InsightsExplorer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var explorers []*models.InsightsExplorer
	err := s.UserDataDB(uid).NewSession(c).Select("explorer_id, uid, name, display_order, hidden").Where("uid=? AND deleted=?", uid, false).Find(&explorers)

	return explorers, err
}

// GetInsightsExplorerByExplorerId returns a insights explorer model according to insights explorer id
func (s *InsightsExplorerService) GetInsightsExplorerByExplorerId(c core.Context, uid int64, explorerId int64) (*models.InsightsExplorer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if explorerId <= 0 {
		return nil, errs.ErrInsightsExplorerIdInvalid
	}

	explorer := &models.InsightsExplorer{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(explorerId).Where("uid=? AND deleted=?", uid, false).Get(explorer)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInsightsExplorerNotFound
	}

	return explorer, nil
}

// GetMaxDisplayOrder returns the max display order
func (s *InsightsExplorerService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	explorer := &models.InsightsExplorer{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(explorer)

	if err != nil {
		return 0, err
	}

	if has {
		return explorer.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreateInsightsExplorer saves a new insights explorer model to database
func (s *InsightsExplorerService) CreateInsightsExplorer(c core.Context, explorer *models.InsightsExplorer) error {
	if explorer.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	explorer.ExplorerId = s.GenerateUuid(uuid.UUID_TYPE_EXPLORER)

	if explorer.ExplorerId < 1 {
		return errs.ErrSystemIsBusy
	}

	explorer.Deleted = false
	explorer.CreatedUnixTime = time.Now().Unix()
	explorer.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(explorer.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(explorer)
		return err
	})
}

// ModifyInsightsExplorer saves an existed insights explorer model to database
func (s *InsightsExplorerService) ModifyInsightsExplorer(c core.Context, explorer *models.InsightsExplorer) error {
	if explorer.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	explorer.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(explorer.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(explorer.ExplorerId).Cols("name", "data", "updated_unix_time").Where("uid=? AND deleted=?", explorer.Uid, false).Update(explorer)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrInsightsExplorerNotFound
		}

		return err
	})
}

// HideInsightsExplorer updates hidden field of given insights explorer ids
func (s *InsightsExplorerService) HideInsightsExplorer(c core.Context, uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.InsightsExplorer{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).In("explorer_id", ids).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrInsightsExplorerNotFound
		}

		return err
	})
}

// ModifyInsightsExplorerDisplayOrders updates display order of given insights explorers
func (s *InsightsExplorerService) ModifyInsightsExplorerDisplayOrders(c core.Context, uid int64, explorers []*models.InsightsExplorer) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(explorers); i++ {
		explorers[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(explorers); i++ {
			explorer := explorers[i]
			updatedRows, err := sess.ID(explorer.ExplorerId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(explorer)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrInsightsExplorerNotFound
			}
		}

		return nil
	})
}

// DeleteInsightsExplorer deletes an existed insights explorer from database
func (s *InsightsExplorerService) DeleteInsightsExplorer(c core.Context, uid int64, explorerId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.InsightsExplorer{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		deletedRows, err := sess.ID(explorerId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrInsightsExplorerNotFound
		}

		return err
	})
}

// DeleteAllInsightsExplorers deletes all existed insights explorers from database
func (s *InsightsExplorerService) DeleteAllInsightsExplorers(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.InsightsExplorer{
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
