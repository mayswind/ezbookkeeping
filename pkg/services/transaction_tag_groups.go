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

// TransactionTagGroupService represents transaction tag group service
type TransactionTagGroupService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a transaction tag group service singleton instance
var (
	TransactionTagGroups = &TransactionTagGroupService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetAllTagGroupsByUid returns all transaction tag group models of user
func (s *TransactionTagGroupService) GetAllTagGroupsByUid(c core.Context, uid int64) ([]*models.TransactionTagGroup, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tagGroups []*models.TransactionTagGroup
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Find(&tagGroups)

	return tagGroups, err
}

// GetTagGroupByTagGroupId returns a transaction tag group model according to transaction tag group id
func (s *TransactionTagGroupService) GetTagGroupByTagGroupId(c core.Context, uid int64, tagGroupId int64) (*models.TransactionTagGroup, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if tagGroupId <= 0 {
		return nil, errs.ErrTransactionTagGroupIdInvalid
	}

	tagGroup := &models.TransactionTagGroup{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(tagGroupId).Where("uid=? AND deleted=?", uid, false).Get(tagGroup)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionTagGroupNotFound
	}

	return tagGroup, nil
}

// GetMaxDisplayOrder returns the max display order
func (s *TransactionTagGroupService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	tagGroup := &models.TransactionTagGroup{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(tagGroup)

	if err != nil {
		return 0, err
	}

	if has {
		return tagGroup.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreateTagGroup saves a new transaction tag group model to database
func (s *TransactionTagGroupService) CreateTagGroup(c core.Context, tagGroup *models.TransactionTagGroup) error {
	if tagGroup.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	tagGroup.TagGroupId = s.GenerateUuid(uuid.UUID_TYPE_TAG_GROUP)

	if tagGroup.TagGroupId < 1 {
		return errs.ErrSystemIsBusy
	}

	tagGroup.Deleted = false
	tagGroup.CreatedUnixTime = time.Now().Unix()
	tagGroup.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tagGroup.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(tagGroup)
		return err
	})
}

// ModifyTagGroup saves an existed transaction tag group model to database
func (s *TransactionTagGroupService) ModifyTagGroup(c core.Context, tagGroup *models.TransactionTagGroup) error {
	if tagGroup.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	tagGroup.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tagGroup.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(tagGroup.TagGroupId).Cols("name", "updated_unix_time").Where("uid=? AND deleted=?", tagGroup.Uid, false).Update(tagGroup)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionTagGroupNotFound
		}

		return err
	})
}

// ModifyTagGroupDisplayOrders updates display order of given transaction tag groups
func (s *TransactionTagGroupService) ModifyTagGroupDisplayOrders(c core.Context, uid int64, tagGroups []*models.TransactionTagGroup) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(tagGroups); i++ {
		tagGroups[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(tagGroups); i++ {
			tagGroup := tagGroups[i]
			updatedRows, err := sess.ID(tagGroup.TagGroupId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(tagGroup)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrTransactionTagGroupNotFound
			}
		}

		return nil
	})
}

// DeleteTagGroup deletes an existed transaction tag group from database
func (s *TransactionTagGroupService) DeleteTagGroup(c core.Context, uid int64, tagGroupId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTagGroup{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		exists, err := sess.Cols("uid", "deleted").Where("uid=? AND deleted=? AND tag_group_id=?", uid, false, tagGroupId).Limit(1).Exist(&models.TransactionTag{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrTransactionTagGroupInUseCannotBeDeleted
		}

		deletedRows, err := sess.ID(tagGroupId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionTagGroupNotFound
		}

		return err
	})
}

// DeleteAllTagGroups deletes all existed transaction tag groups from database
func (s *TransactionTagGroupService) DeleteAllTagGroups(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTagGroup{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		exists, err := sess.Cols("uid", "deleted").Where("uid=? AND deleted=? AND tag_group_id>?", uid, false, 0).Limit(1).Exist(&models.TransactionTag{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrTransactionTagGroupInUseCannotBeDeleted
		}

		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		return nil
	})
}
