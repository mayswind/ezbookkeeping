package api

import (
	"sort"

	"github.com/gin-gonic/gin/binding"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// TransactionCategoriesApi represents transaction category api
type TransactionCategoriesApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	categories *services.TransactionCategoryService
}

// Initialize a transaction category api singleton instance
var (
	TransactionCategories = &TransactionCategoriesApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		categories: services.TransactionCategories,
	}
)

// CategoryListHandler returns transaction category list of current user
func (a *TransactionCategoriesApi) CategoryListHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryListReq models.TransactionCategoryListRequest
	err := c.ShouldBindQuery(&categoryListReq)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	categories, err := a.categories.GetAllCategoriesByUid(c, uid, categoryListReq.Type, categoryListReq.ParentId)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryListHandler] failed to get categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return a.getTransactionCategoryListByTypeResponse(categories, categoryListReq.ParentId)
}

// CategoryGetHandler returns one specific transaction category of current user
func (a *TransactionCategoriesApi) CategoryGetHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryGetReq models.TransactionCategoryGetRequest
	err := c.ShouldBindQuery(&categoryGetReq)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	category, err := a.categories.GetCategoryByCategoryId(c, uid, categoryGetReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryGetHandler] failed to get category \"id:%d\" for user \"uid:%d\", because %s", categoryGetReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	categoryResp := category.ToTransactionCategoryInfoResponse()

	return categoryResp, nil
}

// CategoryCreateHandler saves a new transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryCreateReq models.TransactionCategoryCreateRequest
	err := c.ShouldBindJSON(&categoryCreateReq)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if categoryCreateReq.Type < models.CATEGORY_TYPE_INCOME || categoryCreateReq.Type > models.CATEGORY_TYPE_TRANSFER {
		log.Warnf(c, "[transaction_categories.CategoryCreateHandler] category type invalid, type is %d", categoryCreateReq.Type)
		return nil, errs.ErrTransactionCategoryTypeInvalid
	}

	uid := c.GetCurrentUid()

	if categoryCreateReq.ParentId > 0 {
		parentCategory, err := a.categories.GetCategoryByCategoryId(c, uid, categoryCreateReq.ParentId)

		if err != nil {
			log.Errorf(c, "[transaction_categories.CategoryCreateHandler] failed to get parent category \"id:%d\" for user \"uid:%d\", because %s", categoryCreateReq.ParentId, uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		if parentCategory == nil {
			log.Warnf(c, "[transaction_categories.CategoryCreateHandler] parent category \"id:%d\" does not exist for user \"uid:%d\"", categoryCreateReq.ParentId, uid)
			return nil, errs.ErrParentTransactionCategoryNotFound
		}

		if parentCategory.ParentCategoryId > 0 {
			log.Warnf(c, "[transaction_categories.CategoryCreateHandler] parent category \"id:%d\" has another parent category \"id:%d\" for user \"uid:%d\"", parentCategory.CategoryId, parentCategory.ParentCategoryId, uid)
			return nil, errs.ErrCannotAddToSecondaryTransactionCategory
		}
	}

	var maxOrderId int32

	if categoryCreateReq.ParentId <= 0 {
		maxOrderId, err = a.categories.GetMaxDisplayOrder(c, uid, categoryCreateReq.Type)
	} else {
		maxOrderId, err = a.categories.GetMaxSubCategoryDisplayOrder(c, uid, categoryCreateReq.Type, categoryCreateReq.ParentId)
	}

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	category := a.createNewCategoryModel(uid, &categoryCreateReq, maxOrderId+1)

	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && categoryCreateReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CATEGORY, uid, categoryCreateReq.ClientSessionId)

		if found {
			log.Infof(c, "[transaction_categories.CategoryCreateHandler] another category \"id:%s\" has been created for user \"uid:%d\"", remark, uid)
			categoryId, err := utils.StringToInt64(remark)

			if err == nil {
				category, err = a.categories.GetCategoryByCategoryId(c, uid, categoryId)

				if err != nil {
					log.Errorf(c, "[transaction_categories.CategoryCreateHandler] failed to get existed category \"id:%d\" for user \"uid:%d\", because %s", categoryId, uid, err.Error())
					return nil, errs.Or(err, errs.ErrOperationFailed)
				}

				categoryResp := category.ToTransactionCategoryInfoResponse()

				return categoryResp, nil
			}
		}
	}

	err = a.categories.CreateCategory(c, category)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryCreateHandler] failed to create category \"id:%d\" for user \"uid:%d\", because %s", category.CategoryId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_categories.CategoryCreateHandler] user \"uid:%d\" has created a new category \"id:%d\" successfully", uid, category.CategoryId)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CATEGORY, uid, categoryCreateReq.ClientSessionId, utils.Int64ToString(category.CategoryId))
	categoryResp := category.ToTransactionCategoryInfoResponse()

	return categoryResp, nil
}

// CategoryCreateBatchHandler saves some new transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryCreateBatchHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryCreateBatchReq models.TransactionCategoryCreateBatchRequest
	err := c.ShouldBindBodyWith(&categoryCreateBatchReq, binding.JSON)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryCreateBatchHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	categories, err := a.createBatchCategories(c, uid, &categoryCreateBatchReq)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return a.getTransactionCategoryListByTypeResponse(categories, 0)
}

// CategoryModifyHandler saves an existed transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryModifyReq models.TransactionCategoryModifyRequest
	err := c.ShouldBindJSON(&categoryModifyReq)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	category, err := a.categories.GetCategoryByCategoryId(c, uid, categoryModifyReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryModifyHandler] failed to get category \"id:%d\" for user \"uid:%d\", because %s", categoryModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	newCategory := &models.TransactionCategory{
		CategoryId:       category.CategoryId,
		Uid:              uid,
		ParentCategoryId: categoryModifyReq.ParentId,
		Name:             categoryModifyReq.Name,
		Icon:             categoryModifyReq.Icon,
		Color:            categoryModifyReq.Color,
		Comment:          categoryModifyReq.Comment,
		Hidden:           categoryModifyReq.Hidden,
	}

	if newCategory.ParentCategoryId == category.ParentCategoryId &&
		newCategory.Name == category.Name &&
		newCategory.Icon == category.Icon &&
		newCategory.Color == category.Color &&
		newCategory.Comment == category.Comment &&
		newCategory.Hidden == category.Hidden {
		return nil, errs.ErrNothingWillBeUpdated
	}

	if category.ParentCategoryId == models.LevelOneTransactionCategoryParentId && newCategory.ParentCategoryId != models.LevelOneTransactionCategoryParentId {
		return nil, errs.Or(err, errs.ErrNotAllowChangePrimaryTransactionCategoryToSecondary)
	}

	if category.ParentCategoryId != models.LevelOneTransactionCategoryParentId && newCategory.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
		return nil, errs.Or(err, errs.ErrNotAllowChangeSecondaryTransactionCategoryToPrimary)
	}

	if newCategory.ParentCategoryId != category.ParentCategoryId {
		fromPrimaryCategory, err := a.categories.GetCategoryByCategoryId(c, uid, category.ParentCategoryId)

		if err != nil {
			log.Errorf(c, "[transaction_categories.CategoryModifyHandler] failed to get old primary category \"id:%d\" of category \"id:%d\" for user \"uid:%d\", because %s", category.ParentCategoryId, categoryModifyReq.Id, uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		toPrimaryCategory, err := a.categories.GetCategoryByCategoryId(c, uid, newCategory.ParentCategoryId)

		if err != nil {
			log.Errorf(c, "[transaction_categories.CategoryModifyHandler] failed to get new primary category \"id:%d\" of category \"id:%d\" for user \"uid:%d\", because %s", newCategory.ParentCategoryId, categoryModifyReq.Id, uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		if fromPrimaryCategory.Type != toPrimaryCategory.Type {
			return nil, errs.Or(err, errs.ErrNotAllowChangePrimaryTransactionType)
		}

		if toPrimaryCategory.ParentCategoryId != models.LevelOneTransactionCategoryParentId {
			return nil, errs.Or(err, errs.ErrNotAllowUseSecondaryTransactionAsPrimaryCategory)
		}
	}

	err = a.categories.ModifyCategory(c, newCategory)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryModifyHandler] failed to update category \"id:%d\" for user \"uid:%d\", because %s", categoryModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_categories.CategoryModifyHandler] user \"uid:%d\" has updated category \"id:%d\" successfully", uid, categoryModifyReq.Id)

	newCategory.Type = category.Type
	newCategory.DisplayOrder = category.DisplayOrder
	categoryResp := newCategory.ToTransactionCategoryInfoResponse()

	return categoryResp, nil
}

// CategoryHideHandler hides an existed transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryHideHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryHideReq models.TransactionCategoryHideRequest
	err := c.ShouldBindJSON(&categoryHideReq)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.categories.HideCategory(c, uid, []int64{categoryHideReq.Id}, categoryHideReq.Hidden)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryHideHandler] failed to hide category \"id:%d\" for user \"uid:%d\", because %s", categoryHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_categories.CategoryHideHandler] user \"uid:%d\" has hidden category \"id:%d\"", uid, categoryHideReq.Id)
	return true, nil
}

// CategoryMoveHandler moves display order of existed transaction categories by request parameters for current user
func (a *TransactionCategoriesApi) CategoryMoveHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryMoveReq models.TransactionCategoryMoveRequest
	err := c.ShouldBindJSON(&categoryMoveReq)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryMoveHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	categories := make([]*models.TransactionCategory, len(categoryMoveReq.NewDisplayOrders))

	for i := 0; i < len(categoryMoveReq.NewDisplayOrders); i++ {
		newDisplayOrder := categoryMoveReq.NewDisplayOrders[i]
		category := &models.TransactionCategory{
			Uid:          uid,
			CategoryId:   newDisplayOrder.Id,
			DisplayOrder: newDisplayOrder.DisplayOrder,
		}

		categories[i] = category
	}

	err = a.categories.ModifyCategoryDisplayOrders(c, uid, categories)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryMoveHandler] failed to move categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_categories.CategoryMoveHandler] user \"uid:%d\" has moved categories", uid)
	return true, nil
}

// CategoryDeleteHandler deletes an existed transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var categoryDeleteReq models.TransactionCategoryDeleteRequest
	err := c.ShouldBindJSON(&categoryDeleteReq)

	if err != nil {
		log.Warnf(c, "[transaction_categories.CategoryDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.categories.DeleteCategory(c, uid, categoryDeleteReq.Id)

	if err != nil {
		log.Errorf(c, "[transaction_categories.CategoryDeleteHandler] failed to delete category \"id:%d\" for user \"uid:%d\", because %s", categoryDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_categories.CategoryDeleteHandler] user \"uid:%d\" has deleted category \"id:%d\"", uid, categoryDeleteReq.Id)
	return true, nil
}

func (a *TransactionCategoriesApi) createBatchCategories(c *core.WebContext, uid int64, categoryCreateBatchReq *models.TransactionCategoryCreateBatchRequest) ([]*models.TransactionCategory, error) {
	var err error
	categoryTypeMaxOrderMap := make(map[models.TransactionCategoryType]int32)
	categoriesMap := make(map[*models.TransactionCategory][]*models.TransactionCategory)
	categoriesMap[nil] = make([]*models.TransactionCategory, len(categoryCreateBatchReq.Categories))
	totalCount := 0

	for i := 0; i < len(categoryCreateBatchReq.Categories); i++ {
		categoryCreateReq := categoryCreateBatchReq.Categories[i]
		var maxOrderId, exists = categoryTypeMaxOrderMap[categoryCreateReq.Type]

		if !exists {
			maxOrderId, err = a.categories.GetMaxDisplayOrder(c, uid, categoryCreateReq.Type)

			if err != nil {
				log.Errorf(c, "[transaction_categories.CategoryCreateBatchHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
				return nil, errs.Or(err, errs.ErrOperationFailed)
			}
		}

		category := a.createNewCategoryModel(uid, &models.TransactionCategoryCreateRequest{
			Name:  categoryCreateReq.Name,
			Type:  categoryCreateReq.Type,
			Icon:  categoryCreateReq.Icon,
			Color: categoryCreateReq.Color,
		}, maxOrderId+1)

		categoriesMap[category] = make([]*models.TransactionCategory, len(categoryCreateReq.SubCategories))

		for j := int32(0); j < int32(len(categoryCreateReq.SubCategories)); j++ {
			subCategory := a.createNewCategoryModel(uid, categoryCreateReq.SubCategories[j], j+1)
			categoriesMap[category][j] = subCategory
			totalCount++
		}

		categoriesMap[nil][i] = category
		categoryTypeMaxOrderMap[categoryCreateReq.Type] = maxOrderId + 1
		totalCount++
	}

	categories, err := a.categories.CreateCategories(c, uid, categoriesMap)

	if err != nil {
		log.Errorf(c, "[transaction_categories.createBatchCategories] failed to create categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.Infof(c, "[transaction_categories.createBatchCategories] user \"uid:%d\" has created categories successfully", uid)

	return categories, nil
}

func (a *TransactionCategoriesApi) createNewCategoryModel(uid int64, categoryCreateReq *models.TransactionCategoryCreateRequest, order int32) *models.TransactionCategory {
	return &models.TransactionCategory{
		Uid:              uid,
		Name:             categoryCreateReq.Name,
		Type:             categoryCreateReq.Type,
		ParentCategoryId: categoryCreateReq.ParentId,
		DisplayOrder:     order,
		Icon:             categoryCreateReq.Icon,
		Color:            categoryCreateReq.Color,
		Comment:          categoryCreateReq.Comment,
	}
}

func (a *TransactionCategoriesApi) getTransactionCategoryListByTypeResponse(categories []*models.TransactionCategory, parentId int64) (map[models.TransactionCategoryType]models.TransactionCategoryInfoResponseSlice, *errs.Error) {
	categoryResps := make([]*models.TransactionCategoryInfoResponse, len(categories))
	categoryRespMap := make(map[int64]*models.TransactionCategoryInfoResponse)

	for i := 0; i < len(categories); i++ {
		categoryResps[i] = categories[i].ToTransactionCategoryInfoResponse()
		categoryRespMap[categoryResps[i].Id] = categoryResps[i]
	}

	for i := 0; i < len(categoryResps); i++ {
		categoryResp := categoryResps[i]

		if categoryResp.ParentId <= models.LevelOneTransactionCategoryParentId {
			continue
		}

		parentCategory, parentExists := categoryRespMap[categoryResp.ParentId]

		if !parentExists || parentCategory == nil {
			continue
		}

		parentCategory.SubCategories = append(parentCategory.SubCategories, categoryResp)
	}

	finalCategoryResps := make(models.TransactionCategoryInfoResponseSlice, 0)

	for i := 0; i < len(categoryResps); i++ {
		if parentId <= 0 && categoryResps[i].ParentId == models.LevelOneTransactionCategoryParentId {
			sort.Sort(categoryResps[i].SubCategories)
			finalCategoryResps = append(finalCategoryResps, categoryResps[i])
		} else if parentId > 0 && categoryResps[i].ParentId == parentId {
			finalCategoryResps = append(finalCategoryResps, categoryResps[i])
		}
	}

	sort.Sort(finalCategoryResps)

	typeCategoryMapResponse := make(map[models.TransactionCategoryType]models.TransactionCategoryInfoResponseSlice)

	for i := 0; i < len(finalCategoryResps); i++ {
		category := finalCategoryResps[i]
		categoryList, _ := typeCategoryMapResponse[category.Type]

		categoryList = append(categoryList, category)
		typeCategoryMapResponse[category.Type] = categoryList
	}

	return typeCategoryMapResponse, nil
}
