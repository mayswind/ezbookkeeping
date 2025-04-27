package services

import (
	"strings"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
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

// GetTotalCategoryCountByUid returns total category count of user
func (s *TransactionCategoryService) GetTotalCategoryCountByUid(c core.Context, uid int64) (int64, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	count, err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Count(&models.TransactionCategory{})

	return count, err
}

// GetAllCategoriesByUid returns all transaction category models of user
func (s *TransactionCategoryService) GetAllCategoriesByUid(c core.Context, uid int64, categoryType models.TransactionCategoryType, parentCategoryId int64) ([]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, 8)
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
	err := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("type asc, parent_category_id asc, display_order asc").Find(&categories)

	return categories, err
}

// GetSubCategoriesByCategoryIds returns sub-category models according to category ids
func (s *TransactionCategoryService) GetSubCategoriesByCategoryIds(c core.Context, uid int64, categoryIds []int64) ([]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if len(categoryIds) <= 0 {
		return nil, errs.ErrTransactionCategoryIdInvalid
	}

	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, len(categoryIds)+2)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	var categoryIdConditions strings.Builder

	for i := 0; i < len(categoryIds); i++ {
		if categoryIds[i] <= 0 {
			return nil, errs.ErrTransactionCategoryIdInvalid
		}

		if categoryIdConditions.Len() > 0 {
			categoryIdConditions.WriteString(",")
		}

		categoryIdConditions.WriteString("?")
		conditionParams = append(conditionParams, categoryIds[i])
	}

	if categoryIdConditions.Len() > 1 {
		condition = condition + " AND parent_category_id IN (" + categoryIdConditions.String() + ")"
	} else {
		condition = condition + " AND parent_category_id = " + categoryIdConditions.String()
	}

	var categories []*models.TransactionCategory
	err := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("display_order asc").Find(&categories)

	return categories, err
}

// GetCategoryByCategoryId returns a transaction category model according to transaction category id
func (s *TransactionCategoryService) GetCategoryByCategoryId(c core.Context, uid int64, categoryId int64) (*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if categoryId <= 0 {
		return nil, errs.ErrTransactionCategoryIdInvalid
	}

	category := &models.TransactionCategory{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(categoryId).Where("uid=? AND deleted=?", uid, false).Get(category)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrTransactionCategoryNotFound
	}

	return category, nil
}

// GetCategoriesByCategoryIds returns transaction category models according to transaction category ids
func (s *TransactionCategoryService) GetCategoriesByCategoryIds(c core.Context, uid int64, categoryIds []int64) (map[int64]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if categoryIds == nil {
		return nil, errs.ErrTransactionCategoryIdInvalid
	}

	var categories []*models.TransactionCategory
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).In("category_id", categoryIds).Find(&categories)

	if err != nil {
		return nil, err
	}

	categoryMap := s.GetCategoryMapByList(categories)
	return categoryMap, err
}

// GetMaxDisplayOrder returns the max display order according to transaction category type
func (s *TransactionCategoryService) GetMaxDisplayOrder(c core.Context, uid int64, categoryType models.TransactionCategoryType) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	category := &models.TransactionCategory{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "parent_category_id", "display_order").Where("uid=? AND deleted=? AND type=? AND parent_category_id=?", uid, false, categoryType, models.LevelOneTransactionCategoryParentId).OrderBy("display_order desc").Limit(1).Get(category)

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
func (s *TransactionCategoryService) GetMaxSubCategoryDisplayOrder(c core.Context, uid int64, categoryType models.TransactionCategoryType, parentCategoryId int64) (int32, error) {
	if uid <= 0 {
		return 0, errs.ErrUserIdInvalid
	}

	if parentCategoryId <= 0 {
		return 0, errs.ErrTransactionCategoryIdInvalid
	}

	category := &models.TransactionCategory{}
	has, err := s.UserDataDB(uid).NewSession(c).Cols("uid", "deleted", "parent_category_id", "display_order").Where("uid=? AND deleted=? AND type=? AND parent_category_id=?", uid, false, categoryType, parentCategoryId).OrderBy("display_order desc").Limit(1).Get(category)

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
func (s *TransactionCategoryService) CreateCategory(c core.Context, category *models.TransactionCategory) error {
	if category.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	category.CategoryId = s.GenerateUuid(uuid.UUID_TYPE_CATEGORY)

	if category.CategoryId < 1 {
		return errs.ErrSystemIsBusy
	}

	category.Deleted = false
	category.CreatedUnixTime = time.Now().Unix()
	category.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(category.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(category)
		return err
	})
}

// CreateCategories saves a few transaction category models to database
func (s *TransactionCategoryService) CreateCategories(c core.Context, uid int64, categories map[*models.TransactionCategory][]*models.TransactionCategory) ([]*models.TransactionCategory, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var allCategories []*models.TransactionCategory
	primaryCategories := categories[nil]

	needPrimaryCategoryUuidCount := uint16(len(primaryCategories))
	primaryCategoryUuids := s.GenerateUuids(uuid.UUID_TYPE_CATEGORY, needPrimaryCategoryUuidCount)

	if len(primaryCategoryUuids) < int(needPrimaryCategoryUuidCount) {
		return nil, errs.ErrSystemIsBusy
	}

	for i := 0; i < len(primaryCategories); i++ {
		primaryCategory := primaryCategories[i]
		primaryCategory.CategoryId = primaryCategoryUuids[i]
		primaryCategory.Deleted = false
		primaryCategory.CreatedUnixTime = time.Now().Unix()
		primaryCategory.UpdatedUnixTime = time.Now().Unix()

		allCategories = append(allCategories, primaryCategory)

		secondaryCategories := categories[primaryCategory]

		needSecondaryCategoryUuidCount := uint16(len(secondaryCategories))
		secondaryCategoryUuids := s.GenerateUuids(uuid.UUID_TYPE_CATEGORY, needSecondaryCategoryUuidCount)

		if len(secondaryCategoryUuids) < int(needSecondaryCategoryUuidCount) {
			return nil, errs.ErrSystemIsBusy
		}

		for j := 0; j < len(secondaryCategories); j++ {
			secondaryCategory := secondaryCategories[j]
			secondaryCategory.CategoryId = secondaryCategoryUuids[j]
			secondaryCategory.ParentCategoryId = primaryCategory.CategoryId
			secondaryCategory.Deleted = false
			secondaryCategory.CreatedUnixTime = time.Now().Unix()
			secondaryCategory.UpdatedUnixTime = time.Now().Unix()

			allCategories = append(allCategories, secondaryCategory)
		}
	}

	err := s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionCategoryService) ModifyCategory(c core.Context, category *models.TransactionCategory) error {
	if category.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	category.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(category.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(category.CategoryId).Cols("parent_category_id", "name", "icon", "color", "comment", "hidden", "updated_unix_time").Where("uid=? AND deleted=?", category.Uid, false).Update(category)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrTransactionCategoryNotFound
		}

		return nil
	})
}

// HideCategory updates hidden field of given transaction categories
func (s *TransactionCategoryService) HideCategory(c core.Context, uid int64, ids []int64, hidden bool) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Hidden:          hidden,
		UpdatedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionCategoryService) ModifyCategoryDisplayOrders(c core.Context, uid int64, categories []*models.TransactionCategory) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	for i := 0; i < len(categories); i++ {
		categories[i].UpdatedUnixTime = time.Now().Unix()
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionCategoryService) DeleteCategory(c core.Context, uid int64, categoryId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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
func (s *TransactionCategoryService) DeleteAllCategories(c core.Context, uid int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.TransactionCategory{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
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

// GetVisibleSubCategoryNameMapByList returns visible sub transaction category map by a list
func (s *TransactionCategoryService) GetVisibleSubCategoryNameMapByList(categories []*models.TransactionCategory) (expenseCategoryMap map[string]map[string]*models.TransactionCategory, incomeCategoryMap map[string]map[string]*models.TransactionCategory, transferCategoryMap map[string]map[string]*models.TransactionCategory) {
	categoryMap := make(map[int64]*models.TransactionCategory, len(categories))
	expenseCategoryMap = make(map[string]map[string]*models.TransactionCategory)
	incomeCategoryMap = make(map[string]map[string]*models.TransactionCategory)
	transferCategoryMap = make(map[string]map[string]*models.TransactionCategory)

	for i := 0; i < len(categories); i++ {
		category := categories[i]
		categoryMap[category.CategoryId] = category
	}

	for i := 0; i < len(categories); i++ {
		category := categories[i]

		if category.Hidden {
			continue
		}

		if category.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
			continue
		}

		parentCategory, exists := categoryMap[category.ParentCategoryId]

		if !exists {
			continue
		}

		var categories map[string]*models.TransactionCategory

		if category.Type == models.CATEGORY_TYPE_INCOME {
			categories, exists = incomeCategoryMap[category.Name]

			if !exists {
				categories = make(map[string]*models.TransactionCategory)
				incomeCategoryMap[category.Name] = categories
			}
		} else if category.Type == models.CATEGORY_TYPE_EXPENSE {
			categories, exists = expenseCategoryMap[category.Name]

			if !exists {
				categories = make(map[string]*models.TransactionCategory)
				expenseCategoryMap[category.Name] = categories
			}
		} else if category.Type == models.CATEGORY_TYPE_TRANSFER {
			categories, exists = transferCategoryMap[category.Name]

			if !exists {
				categories = make(map[string]*models.TransactionCategory)
				transferCategoryMap[category.Name] = categories
			}
		} else {
			continue
		}

		categories[parentCategory.Name] = category
	}

	return expenseCategoryMap, incomeCategoryMap, transferCategoryMap
}

// GetCategoryNames returns a list with transaction category names from transaction category models list
func (s *TransactionCategoryService) GetCategoryNames(categories []*models.TransactionCategory) []string {
	categoryNames := make([]string, len(categories))

	for i := 0; i < len(categories); i++ {
		categoryNames[i] = categories[i].Name
	}

	return categoryNames
}
