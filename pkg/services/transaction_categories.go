package services

import (
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/lab/pkg/datastore"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/uuid"
)

type TransactionCategoryService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

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

func (s *TransactionCategoryService) GetAllCategoriesByUid(uid int64, categoryType models.TransactionCategoryType, parentCategoryId int64) ([]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	condition := "uid=? AND deleted=?"
	var conditionParams []interface{}
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

func (s *TransactionCategoryService) GetCategoryAndSubCategoriesByCategoryId(uid int64, categoryId int64) ([]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if categoryId <= 0 {
		return nil, errs.ErrTransactionCategoryIdInvalid
	}

	var categories []*models.TransactionCategory
	err := s.UserDataDB(uid).Where("uid=? AND deleted=? AND (category_id=? OR parent_category_id=?)", uid, false, categoryId, categoryId).OrderBy("type asc, parent_category_id asc, display_order asc").Find(&categories)

	return categories, err
}

func (s *TransactionCategoryService) GetMaxDisplayOrder(uid int64, categoryType models.TransactionCategoryType) (int, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	category := &models.TransactionCategory{}
	has, err := s.UserDataDB(uid).Cols("uid", "deleted", "parent_category_id", "display_order").Where("uid=? AND deleted=? AND type=? AND parent_category_id=?", uid, false, categoryType, models.TRANSACTION_PARENT_ID_LEVEL_ONE).OrderBy("display_order desc").Limit(1).Get(category)

	if err != nil {
		return 0, err
	}

	if has {
		return category.DisplayOrder, nil
	} else {
		return 0, nil
	}
}

func (s *TransactionCategoryService) GetMaxSubCategoryDisplayOrder(uid int64, categoryType models.TransactionCategoryType, parentCategoryId int64) (int, error) {
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

func (s *TransactionCategoryService) HideCategory(uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Hidden: hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		updatedRows, err := sess.Cols("hidden", "updated_unix_time").In("category_id", ids).Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionCategoryNotFound
		}

		return nil
	})
}

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

func (s *TransactionCategoryService) DeleteCategories(uid int64, ids []int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Deleted: true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(func(sess *xorm.Session) error {
		deletedRows, err := sess.Cols("deleted", "deleted_unix_time").In("category_id", ids).Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrTransactionCategoryNotFound
		}

		_, err = sess.Cols("deleted", "deleted_unix_time").In("parent_category_id", ids).Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		return err
	})
}
