package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// TransactionCategoryService represents transaction category service
type TransactionCategoryService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a transaction category service singleton instance
var (
	TransactionCategories = &TransactionCategoryService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetAllCategoriesByUid returns all transaction category models of user
func (s *TransactionCategoryService) GetAllCategoriesByUid(uid int64, categoryType models.TransactionCategoryType, parentCategoryId int64) ([]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	condition := "uid=? AND deleted=?"
	conditionParams := make([]interface{}, 0, 8)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	if categoryType > 0 {
		condition = condition + " AND type=?"
		conditionParams = append(conditionParams, categoryType)
	}

	if parentCategoryId >= 0 {
		condition = condition + " AND parent_category_id=?"
		conditionParams = append(conditionParams, parentCategoryId)
	}

	var categories []*models.TransactionCategory
	err := s.UserDataDB(uid).Where(condition, conditionParams...).OrderBy("type asc, parent_category_id asc, display_order asc").Find(&categories)

	return categories, err
}

// GetCategoryByCategoryId returns a transaction category model according to transaction category id
func (s *TransactionCategoryService) GetCategoryByCategoryId(uid int64, categoryId int64) (*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if categoryId <= 0 {
		return nil, errs.ErrTransactionCategoryIdInvalid
	}

	category := &models.TransactionCategory{}
	has, err := s.UserDataDB(uid).ID(categoryId).Where("uid=? AND deleted=?", uid, false).Get(category)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionCategoryNotFound
	}

	return category, nil
}

// GetCategoriesByCategoryIds returns transaction category models according to transaction category ids
func (s *TransactionCategoryService) GetCategoriesByCategoryIds(uid int64, categoryIds []int64) (map[int64]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if categoryIds == nil {
		return nil, errs.ErrTransactionCategoryIdInvalid
	}

	var categories []*models.TransactionCategory
	err := s.UserDataDB(uid).Where("uid=? AND deleted=?", uid, false).In("category_id", categoryIds).Find(&categories)

	if err != nil {
		return nil, err
	}

	categoryMap := s.GetCategoryMapByList(categories)
	return categoryMap, err
}

// GetMaxDisplayOrder returns the max display order according to transaction category type
func (s *TransactionCategoryService) GetMaxDisplayOrder(uid int64, categoryType models.TransactionCategoryType) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	category := &models.TransactionCategory{}
	has, err := s.UserDataDB(uid).Cols("uid", "deleted", "parent_category_id", "display_order").Where("uid=? AND deleted=? AND type=? AND parent_category_id=?", uid, false, categoryType, models.LevelOneTransactionParentId).OrderBy("display_order desc").Limit(1).Get(category)

	if err != nil {
		return 0, err
	}

	if has {
		return category.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// GetMaxSubCategoryDisplayOrder returns the max display order of sub transaction category according to transaction category type and parent transaction category id
func (s *TransactionCategoryService) GetMaxSubCategoryDisplayOrder(uid int64, categoryType models.TransactionCategoryType, parentCategoryId int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	if parentCategoryId <= 0 {
		return 0, errs.ErrTransactionCategoryIdInvalid
	}

	category := &models.TransactionCategory{}
	has, err := s.UserDataDB(uid).Cols("uid", "deleted", "parent_category_id", "display_order").Where("uid=? AND deleted=? AND type=? AND parent_category_id=?", uid, false, categoryType, parentCategoryId).OrderBy("display_order desc").Limit(1).Get(category)

	if err != nil {
		return 0, err
	}

	if has {
		return category.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

// CreateCategory saves a new transaction category model to database
func (s *TransactionCategoryService) CreateCategory(category *models.TransactionCategory) error {
	if category.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	category.CategoryId = s.GenerateUuid(uuid.UUID_TYPE_CATEGORY)

	category.Deleted = false
	category.CreatedUnixTime = time.Now().Unix()
	category.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(category.Uid).DoTransaction(func(sess *xorm.Session) error {
		_, err := sess.Insert(category)
		return err
	})
}

// CreateCategories saves a few transaction category models to database
func (s *TransactionCategoryService) CreateCategories(uid int64, categories map[*models.TransactionCategory][]*models.TransactionCategory) ([]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var allCategories []*models.TransactionCategory
	primaryCategories := categories[nil]

	for i := 0; i < len(primaryCategories); i++ {
		primaryCategory := primaryCategories[i]
		primaryCategory.CategoryId = s.GenerateUuid(uuid.UUID_TYPE_CATEGORY)

		primaryCategory.Deleted = false
		primaryCategory.CreatedUnixTime = time.Now().Unix()
		primaryCategory.UpdatedUnixTime = time.Now().Unix()

		allCategories = append(allCategories, primaryCategory)

		secondaryCategories := categories[primaryCategory]

		for j := 0; j < len(secondaryCategories); j++ {
			secondaryCategory := secondaryCategories[j]
			secondaryCategory.CategoryId = s.GenerateUuid(uuid.UUID_TYPE_CATEGORY)
			secondaryCategory.ParentCategoryId = primaryCategory.CategoryId

			secondaryCategory.Deleted = false
			secondaryCategory.CreatedUnixTime = time.Now().Unix()
			secondaryCategory.UpdatedUnixTime = time.Now().Unix()

			allCategories = append(allCategories, secondaryCategory)
		}
	}

	err := s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(allCategories); i++ {
			category := allCategories[i]
			_, err := sess.Insert(category)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return allCategories, nil
}

// ModifyCategory saves an existed transaction category model to database
func (s *TransactionCategoryService) ModifyCategory(category *models.TransactionCategory) error {
	if category.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	category.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(category.Uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(category.CategoryId).Cols("name", "icon", "color", "comment", "hidden", "updated_unix_time").Where("uid=? AND deleted=?", category.Uid, false).Update(category)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionCategoryNotFound
		}

		return nil
	})
}

// HideCategory updates hidden field of given transaction categories
func (s *TransactionCategoryService) HideCategory(uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).In("category_id", ids).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionCategoryNotFound
		}

		return nil
	})
}

// ModifyCategoryDisplayOrders updates display order of given transaction categories
func (s *TransactionCategoryService) ModifyCategoryDisplayOrders(uid int64, categories []*models.TransactionCategory) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(categories); i++ {
		categories[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		for i := 0; i < len(categories); i++ {
			category := categories[i]
			updatedRows, err := sess.ID(category.CategoryId).Cols("display_order", "updated_unix_time").Where("uid=? AND deleted=?", uid, false).Update(category)

			if err != nil {
				return err
			} else if updatedRows < 1 {
				return errs.ErrTransactionCategoryNotFound
			}
		}

		return nil
	})
}

// DeleteCategory deletes an existed transaction category from database
func (s *TransactionCategoryService) DeleteCategory(uid int64, categoryId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		var categoryAndSubCategories []*models.TransactionCategory
		err := sess.Where("uid=? AND deleted=? AND (category_id=? OR parent_category_id=?)", uid, false, categoryId, categoryId).Find(&categoryAndSubCategories)

		if err != nil {
			return err
		} else if len(categoryAndSubCategories) < 1 {
			return errs.ErrTransactionCategoryNotFound
		}

		categoryAndSubCategoryIds := make([]int64, len(categoryAndSubCategories))

		for i := 0; i < len(categoryAndSubCategories); i++ {
			categoryAndSubCategoryIds[i] = categoryAndSubCategories[i].CategoryId
		}

		exists, err := sess.Cols("uid", "deleted", "category_id").Where("uid=? AND deleted=?", uid, false).In("category_id", categoryAndSubCategoryIds).Limit(1).Exist(&models.Transaction{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrTransactionCategoryInUseCannotBeDeleted
		}

		deletedRows, err := sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).In("category_id", categoryAndSubCategoryIds).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionCategoryNotFound
		}

		return err
	})
}

// DeleteAllCategories deletes all existed transaction categories from database
func (s *TransactionCategoryService) DeleteAllCategories(uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		exists, err := sess.Cols("uid", "deleted", "category_id").Where("uid=? AND deleted=? AND category_id<>?", uid, false, 0).Limit(1).Exist(&models.Transaction{})

		if err != nil {
			return err
		} else if exists {
			return errs.ErrTransactionCategoryInUseCannotBeDeleted
		}

		_, err = sess.Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		}

		return nil
	})
}

// GetCategoryMapByList returns a transaction category map by a list
func (s *TransactionCategoryService) GetCategoryMapByList(categories []*models.TransactionCategory) map[int64]*models.TransactionCategory {
	categoryMap := make(map[int64]*models.TransactionCategory)

	for i := 0; i < len(categories); i++ {
		category := categories[i]
		categoryMap[category.CategoryId] = category
	}
	return categoryMap
}
