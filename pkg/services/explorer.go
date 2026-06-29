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

// GetTotalExplorationsCountByUid returns total explorations count of user
func (s *InsightsExplorerService) GetTotalExplorationsCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.InsightsExplorer{})

	return count, err
}

// GetAllExplorationNamesByUid returns all exploration names
func (s *InsightsExplorerService) GetAllExplorationNamesByUid(c core.Context, uid int64) ([]*models.InsightsExplorer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var explorers []*models.InsightsExplorer
	err := s.UserDataDB(uid).NewSession(c).Select("explorer_id, uid, name, display_order, hidden").Where("uid=? AND deleted=?", uid, false).Find(&explorers)

	return explorers, err
}

// GetExplorationByExplorationId returns a exploration model according to exploration id
func (s *InsightsExplorerService) GetExplorationByExplorationId(c core.Context, uid int64, explorationId int64) (*models.InsightsExplorer, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if explorationId <= 0 {
		return nil, errs.ErrInsightsExplorerIdInvalid
	}

	exploration := &models.InsightsExplorer{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(explorationId).Where("uid=? AND deleted=?", uid, false).Get(exploration)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInsightsExplorerNotFound
	}

	return exploration, nil
}

// GetMaxDisplayOrder returns the max display order
func (s *InsightsExplorerService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	exploration := &models.InsightsExplorer{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(exploration)

	if err != nil {
		return 0, err
	}

	if has {
		return exploration.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreateExploration saves a new exploration model to database
func (s *InsightsExplorerService) CreateExploration(c core.Context, exploration *models.InsightsExplorer) error {
	if exploration.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exploration.ExplorerId = s.GenerateUuid(uuid.UUID_TYPE_EXPLORER)

	if exploration.ExplorerId < 1 {
		return errs.ErrSystemIsBusy
	}

	exploration.Deleted = false
	exploration.CreatedUnixTime = time.Now().Unix()
	exploration.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(exploration.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(exploration)
		return err
	})
}

// ModifyExploration saves an existed exploration model to database
func (s *InsightsExplorerService) ModifyExploration(c core.Context, exploration *models.InsightsExplorer) error {
	if exploration.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exploration.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(exploration.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(exploration.ExplorerId).Cols("name", "data", "updated_unix_time").Where("uid=? AND deleted=?", exploration.Uid, false).Update(exploration)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrInsightsExplorerNotFound
		}

		return err
	})
}

// HideExploration updates hidden field of given exploration ids
func (s *InsightsExplorerService) HideExploration(c core.Context, uid int64, ids []int64, hidden bool) error {
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

// ModifyExplorationDisplayOrders updates display order of given explorations
func (s *InsightsExplorerService) ModifyExplorationDisplayOrders(c core.Context, uid int64, explorations []*models.InsightsExplorer) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(explorations); i++ {
		explorations[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(explorations); i++ {
			exploration := explorations[i]
			updatedRows, err := sess.ID(exploration.ExplorerId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(exploration)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrInsightsExplorerNotFound
			}
		}

		return nil
	})
}

// DeleteExploration deletes an existed exploration from database
func (s *InsightsExplorerService) DeleteExploration(c core.Context, uid int64, explorerId int64) error {
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

// DeleteAllExplorations deletes all existed explorations from database
func (s *InsightsExplorerService) DeleteAllExplorations(c core.Context, uid int64) error {
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
