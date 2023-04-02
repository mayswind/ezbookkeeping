package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// TransactionTagService represents transaction tag service
type TransactionTagService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a transaction tag service singleton instance
var (
	TransactionTags = &TransactionTagService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetTotalTagCountByUid returns total tag count of user
func (s *TransactionTagService) GetTotalTagCountByUid(uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).Count(&models.TransactionTag{})

	return count, err
}

// GetAllTagsByUid returns all transaction tag models of user
func (s *TransactionTagService) GetAllTagsByUid(uid int64) ([]*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tags []*models.TransactionTag
	err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).Find(&tags)

	return tags, err
}

// GetTagByTagId returns a transaction tag model according to transaction tag id
func (s *TransactionTagService) GetTagByTagId(uid int64, tagId int64) (*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if tagId <= 0 {
		return nil, errs.ErrTransactionTagIdInvalid
	}

	tag := &models.TransactionTag{}
	has, err := s.UserDataDB(uid).ID(tagId).Where("uid=? AND deleted=?", uid, false).Get(tag)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionTagNotFound
	}

	return tag, nil
}

// GetTagsByTagIds returns transaction tag models according to transaction tag ids
func (s *TransactionTagService) GetTagsByTagIds(uid int64, tagIds []int64) (map[int64]*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if tagIds == nil {
		return nil, errs.ErrTransactionTagIdInvalid
	}

	var tags []*models.TransactionTag
	err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).In("tag_id", tagIds).Find(&tags)

	if err != nil {
		return nil, err
	}

	tagMap := s.GetTagMapByList(tags)
	return tagMap, err
}

// GetMaxDisplayOrder returns the max display order
func (s *TransactionTagService) GetMaxDisplayOrder(uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	tag := &models.TransactionTag{}
	has, err := s.UserDataDB(uid).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(tag)

	if err != nil {
		return 0, err
	}

	if has {
		return tag.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// GetAllTagIdsOfAllTransactions returns all transaction tag ids
func (s *TransactionTagService) GetAllTagIdsOfAllTransactions(uid int64) (map[int64][]int64, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tagIndexs []*models.TransactionTagIndex
	err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).Find(&tagIndexs)

	allTransactionTagIds := s.getGroupedTransactionTagIds(tagIndexs)

	return allTransactionTagIds, err
}

// GetAllTagIdsOfTransactions returns transaction tag ids for given transactions
func (s *TransactionTagService) GetAllTagIdsOfTransactions(uid int64, transactionIds []int64) (map[int64][]int64, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tagIndexs []*models.TransactionTagIndex
	err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).In("transaction_id", transactionIds).Find(&tagIndexs)

	allTransactionTagIds := s.getGroupedTransactionTagIds(tagIndexs)

	return allTransactionTagIds, err
}

// CreateTag saves a new transaction tag model to database
func (s *TransactionTagService) CreateTag(tag *models.TransactionTag) error {
	if tag.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsTagName(tag.Uid, tag.Name)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrTransactionTagNameAlreadyExists
	}

	tag.TagId = s.GenerateUuid(uuid.UUID_TYPE_TAG)

	tag.Deleted = false
	tag.CreatedUnixTime = time.Now().Unix()
	tag.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tag.Uid).DoTransaction(func(sess *xorm.Session) error {
		_, err := sess.Insert(tag)
		return err
	})
}

// ModifyTag saves an existed transaction tag model to database
func (s *TransactionTagService) ModifyTag(tag *models.TransactionTag) error {
	if tag.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsTagName(tag.Uid, tag.Name)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrTransactionTagNameAlreadyExists
	}

	tag.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tag.Uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(tag.TagId).Cols("name", "updated_unix_time").Where("uid=? AND deleted=?", tag.Uid, false).Update(tag)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionTagNotFound
		}

		return err
	})
}

// HideTag updates hidden field of given transaction tags
func (s *TransactionTagService) HideTag(uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTag{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).In("tag_id", ids).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionTagNotFound
		}

		return err
	})
}

// ModifyTagDisplayOrders updates display order of given transaction tags
func (s *TransactionTagService) ModifyTagDisplayOrders(uid int64, tags []*models.TransactionTag) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(tags); i++ {
		tags[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(tags); i++ {
			tag := tags[i]
			updatedRows, err := sess.ID(tag.TagId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(tag)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrTransactionTagNotFound
			}
		}

		return nil
	})
}

// DeleteTag deletes an existed transaction tag from database
func (s *TransactionTagService) DeleteTag(uid int64, tagId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTag{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		exists, err := sess.Cols("uid", "tag_id").Where("uid=? AND deleted=? AND tag_id=?", uid, false, tagId).Limit(1).Exist(&models.TransactionTagIndex{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrTransactionTagInUseCannotBeDeleted
		}

		deletedRows, err := sess.ID(tagId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionTagNotFound
		}

		return err
	})
}

// DeleteAllTags deletes all existed transaction tags from database
func (s *TransactionTagService) DeleteAllTags(uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTag{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		exists, err := sess.Cols("uid", "deleted").Where("uid=? AND deleted=?", uid, false).Limit(1).Exist(&models.TransactionTagIndex{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrTransactionTagInUseCannotBeDeleted
		}

		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		return nil
	})
}

// ExistsTagName returns whether the given tag name exists
func (s *TransactionTagService) ExistsTagName(uid int64, name string) (bool, error) {
	if name == "" {
		return false, errs.ErrTransactionTagNameIsEmpty
	}

	return s.UserDataDB(uid).Cols("name").Where("uid=? AND deleted=? AND name=?", uid, false, name).Exist(&models.TransactionTag{})
}

// GetTagMapByList returns a transaction tag map by a list
func (s *TransactionTagService) GetTagMapByList(tags []*models.TransactionTag) map[int64]*models.TransactionTag {
	tagMap := make(map[int64]*models.TransactionTag)

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		tagMap[tag.TagId] = tag
	}
	return tagMap
}

func (s *TransactionTagService) getGroupedTransactionTagIds(tagIndexs []*models.TransactionTagIndex) map[int64][]int64 {
	allTransactionTagIds := make(map[int64][]int64)

	for i := 0; i < len(tagIndexs); i++ {
		tagIndex := tagIndexs[i]

		var transactionTagIds []int64

		if _, exists := allTransactionTagIds[tagIndex.TransactionId]; exists {
			transactionTagIds = allTransactionTagIds[tagIndex.TransactionId]
		}

		transactionTagIds = append(transactionTagIds, tagIndex.TagId)
		allTransactionTagIds[tagIndex.TransactionId] = transactionTagIds
	}
	return allTransactionTagIds
}
