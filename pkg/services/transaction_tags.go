package services

import (
	"strings"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

const pageCountForLoadAllTransactionTagIndexes = 1000

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
func (s *TransactionTagService) GetTotalTagCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.TransactionTag{})

	return count, err
}

// GetAllTagsByUid returns all transaction tag models of user
func (s *TransactionTagService) GetAllTagsByUid(c core.Context, uid int64) ([]*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tags []*models.TransactionTag
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Find(&tags)

	return tags, err
}

// GetTagByTagId returns a transaction tag model according to transaction tag id
func (s *TransactionTagService) GetTagByTagId(c core.Context, uid int64, tagId int64) (*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if tagId <= 0 {
		return nil, errs.ErrTransactionTagIdInvalid
	}

	tag := &models.TransactionTag{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(tagId).Where("uid=? AND deleted=?", uid, false).Get(tag)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionTagNotFound
	}

	return tag, nil
}

// GetTagsByTagIds returns transaction tag models according to transaction tag ids
func (s *TransactionTagService) GetTagsByTagIds(c core.Context, uid int64, tagIds []int64) (map[int64]*models.TransactionTag, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if tagIds == nil {
		return nil, errs.ErrTransactionTagIdInvalid
	}

	var tags []*models.TransactionTag
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("tag_id", tagIds).Find(&tags)

	if err != nil {
		return nil, err
	}

	tagMap := s.GetTagMapByList(tags)
	return tagMap, err
}

// GetMaxDisplayOrder returns the max display order
func (s *TransactionTagService) GetMaxDisplayOrder(c core.Context, uid int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	tag := &models.TransactionTag{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "display_order").Where("uid=? AND deleted=?", uid, false).OrderBy("display_order desc").Limit(1).Get(tag)

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
func (s *TransactionTagService) GetAllTagIdsOfAllTransactions(c core.Context, uid int64) ([]*models.TransactionTagIndex, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, 2)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	var allTransactionTagIndexes []*models.TransactionTagIndex

	maxTransactionTagIndexId := int64(0)

	for maxTransactionTagIndexId >= 0 {
		var tagIndexes []*models.TransactionTagIndex

		finalCondition := condition
		finalConditionParams := make([]any, 0, 3)
		finalConditionParams = append(finalConditionParams, conditionParams...)

		if maxTransactionTagIndexId > 0 {
			finalCondition = finalCondition + " AND tag_index_id<=?"
			finalConditionParams = append(finalConditionParams, maxTransactionTagIndexId)
		}

		err := s.UserDataDB(uid).NewSession(c).Where(finalCondition, finalConditionParams...).Limit(pageCountForLoadAllTransactionTagIndexes, 0).OrderBy("tag_index_id desc").Find(&tagIndexes)

		if err != nil {
			return nil, err
		}

		allTransactionTagIndexes = append(allTransactionTagIndexes, tagIndexes...)

		if len(tagIndexes) < pageCountForLoadAllTransactionTagIndexes {
			maxTransactionTagIndexId = -1
			break
		}

		maxTransactionTagIndexId = tagIndexes[len(tagIndexes)-1].TagIndexId - 1
	}

	return allTransactionTagIndexes, nil
}

// GetAllTagIdsMapOfAllTransactions returns all transaction tag ids map grouped by transaction id
func (s *TransactionTagService) GetAllTagIdsMapOfAllTransactions(c core.Context, uid int64) (map[int64][]int64, error) {
	tagIndexes, err := s.GetAllTagIdsOfAllTransactions(c, uid)

	if err != nil {
		return nil, err
	}

	allTransactionTagIds := s.GetGroupedTransactionTagIds(tagIndexes)

	return allTransactionTagIds, err
}

// GetAllTagIdsOfTransactions returns transaction tag ids for given transactions
func (s *TransactionTagService) GetAllTagIdsOfTransactions(c core.Context, uid int64, transactionIds []int64) (map[int64][]int64, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var tagIndexes []*models.TransactionTagIndex
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("transaction_id", transactionIds).OrderBy("transaction_id asc, tag_index_id asc").Find(&tagIndexes)

	allTransactionTagIds := s.GetGroupedTransactionTagIds(tagIndexes)

	return allTransactionTagIds, err
}

// CreateTag saves a new transaction tag model to database
func (s *TransactionTagService) CreateTag(c core.Context, tag *models.TransactionTag) error {
	if tag.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsTagName(c, tag.Uid, tag.Name)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrTransactionTagNameAlreadyExists
	}

	tag.TagId = s.GenerateUuid(uuid.UUID_TYPE_TAG)

	if tag.TagId < 1 {
		return errs.ErrSystemIsBusy
	}

	tag.Deleted = false
	tag.CreatedUnixTime = time.Now().Unix()
	tag.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tag.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(tag)
		return err
	})
}

// CreateTags saves a few transaction tag models to database
func (s *TransactionTagService) CreateTags(c core.Context, uid int64, tags []*models.TransactionTag, skipExists bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	allTagNames := make([]string, len(tags))

	for i := 0; i < len(tags); i++ {
		allTagNames[i] = tags[i].Name
	}

	var existTags []*models.TransactionTag
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("name", allTagNames).Find(&existTags)

	if err != nil {
		return err
	} else if !skipExists && len(existTags) > 0 {
		return errs.ErrTransactionTagNameAlreadyExists
	}

	existsNameTagMap := make(map[string]*models.TransactionTag, len(existTags))

	for i := 0; i < len(existTags); i++ {
		tag := existTags[i]
		existsNameTagMap[tag.Name] = tag
	}

	newTags := make([]*models.TransactionTag, 0, len(tags)-len(existTags))

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		existsTag, exists := existsNameTagMap[tag.Name]

		if exists {
			tag.FillFromOtherTag(existsTag)
			continue
		}

		newTags = append(newTags, tag)
	}

	tagUuids := s.GenerateUuids(uuid.UUID_TYPE_TAG_INDEX, uint16(len(newTags)))

	if len(tagUuids) < len(newTags) {
		return errs.ErrSystemIsBusy
	}

	for i := 0; i < len(newTags); i++ {
		tag := newTags[i]
		tag.TagId = tagUuids[i]
		tag.Deleted = false
		tag.CreatedUnixTime = time.Now().Unix()
		tag.UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(newTags); i++ {
			tag := newTags[i]
			_, err := sess.Insert(tag)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

// ModifyTag saves an existed transaction tag model to database
func (s *TransactionTagService) ModifyTag(c core.Context, tag *models.TransactionTag) error {
	if tag.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	exists, err := s.ExistsTagName(c, tag.Uid, tag.Name)

	if err != nil {
		return err
	} else if exists {
		return errs.ErrTransactionTagNameAlreadyExists
	}

	tag.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(tag.Uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionTagService) HideTag(c core.Context, uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTag{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionTagService) ModifyTagDisplayOrders(c core.Context, uid int64, tags []*models.TransactionTag) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(tags); i++ {
		tags[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionTagService) DeleteTag(c core.Context, uid int64, tagId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTag{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		exists, err := sess.Cols("uid", "tag_id").Where("uid=? AND deleted=? AND tag_id=?", uid, false, tagId).Limit(1).Exist(&models.TransactionTagIndex{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrTransactionTagInUseCannotBeDeleted
		}

		var relatedTransactionTemplatesByTag []*models.TransactionTemplate
		err = sess.Cols("uid", "deleted", "tag_ids", "template_type", "scheduled_frequency_type", "scheduled_end_time").Where("uid=? AND deleted=? AND (template_type=? OR (template_type=? AND scheduled_frequency_type<>? AND (scheduled_end_time IS NULL OR scheduled_end_time>=?))) AND tag_ids LIKE ?", uid, false, models.TRANSACTION_TEMPLATE_TYPE_NORMAL, models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE, models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED, now, "%%"+utils.Int64ToString(tagId)+"%%").Find(&relatedTransactionTemplatesByTag)

		if err != nil {
			return err
		}

		for i := 0; i < len(relatedTransactionTemplatesByTag); i++ {
			template := relatedTransactionTemplatesByTag[i]
			tagIds, err := s.GetTagIds(template.TagIds)

			if err != nil {
				return err
			}

			for j := 0; j < len(tagIds); j++ {
				if tagIds[j] == tagId {
					return errs.ErrTransactionTagInUseCannotBeDeleted
				}
			}
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
func (s *TransactionTagService) DeleteAllTags(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionTag{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionTagService) ExistsTagName(c core.Context, uid int64, name string) (bool, error) {
	if name == "" {
		return false, errs.ErrTransactionTagNameIsEmpty
	}

	return s.UserDataDB(uid).NewSession(c).Cols("name").Where("uid=? AND deleted=? AND name=?", uid, false, name).Exist(&models.TransactionTag{})
}

// ModifyTagIndexTransactionTime updates transaction time of given transaction tag indexes
func (s *TransactionTagService) ModifyTagIndexTransactionTime(c core.Context, uid int64, tagIndexes []*models.TransactionTagIndex) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(tagIndexes); i++ {
		tagIndexes[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		for i := 0; i < len(tagIndexes); i++ {
			tagIndex := tagIndexes[i]
			updatedRows, err := sess.ID(tagIndex.TagIndexId).Cols("transaction_time", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(tagIndex)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrTransactionTagIndexNotFound
			}
		}

		return nil
	})
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

// GetTagNameMapByList returns a transaction tag map by a list
func (s *TransactionTagService) GetTagNameMapByList(tags []*models.TransactionTag) map[string]*models.TransactionTag {
	tagMap := make(map[string]*models.TransactionTag)

	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		tagMap[tag.Name] = tag
	}
	return tagMap
}

// GetTagNames returns a list with tag names from tag models list
func (s *TransactionTagService) GetTagNames(tags []*models.TransactionTag) []string {
	tagNames := make([]string, len(tags))

	for i := 0; i < len(tags); i++ {
		tagNames[i] = tags[i].Name
	}

	return tagNames
}

// GetGroupedTransactionTagIds returns a map of transaction tag ids grouped by transaction id
func (s *TransactionTagService) GetGroupedTransactionTagIds(tagIndexes []*models.TransactionTagIndex) map[int64][]int64 {
	allTransactionTagIds := make(map[int64][]int64)

	for i := 0; i < len(tagIndexes); i++ {
		tagIndex := tagIndexes[i]

		var transactionTagIds []int64

		if _, exists := allTransactionTagIds[tagIndex.TransactionId]; exists {
			transactionTagIds = allTransactionTagIds[tagIndex.TransactionId]
		}

		transactionTagIds = append(transactionTagIds, tagIndex.TagId)
		allTransactionTagIds[tagIndex.TransactionId] = transactionTagIds
	}

	for _, tagIds := range allTransactionTagIds {
		utils.Int64Sort(tagIds)
	}

	return allTransactionTagIds
}

// GetTagIds converts a comma-separated string of tag ids into a slice of int64
func (s *TransactionTagService) GetTagIds(tagIds string) ([]int64, error) {
	if tagIds == "" || tagIds == "0" {
		return nil, nil
	}

	requestTagIds, err := utils.StringArrayToInt64Array(strings.Split(tagIds, ","))

	if err != nil {
		return nil, errs.Or(err, errs.ErrTransactionTagIdInvalid)
	}

	return requestTagIds, nil
}

// GetTransactionTagIds returns a slice of all transaction tag ids from a map of transaction tag ids grouped by transaction id
func (s *TransactionTagService) GetTransactionTagIds(allTransactionTagIds map[int64][]int64) []int64 {
	allTagIds := make([]int64, 0, len(allTransactionTagIds))

	for _, tagIds := range allTransactionTagIds {
		allTagIds = append(allTagIds, tagIds...)
	}

	return allTagIds
}
