package api

import (
	"sort"

	"github.com/mayswind/lab/pkg/core"
	"github.com/mayswind/lab/pkg/errs"
	"github.com/mayswind/lab/pkg/log"
	"github.com/mayswind/lab/pkg/models"
	"github.com/mayswind/lab/pkg/services"
)

// TransactionCategoriesApi represents transaction category api
type TransactionCategoriesApi struct {
	categories *services.TransactionCategoryService
}

// Initialize a transaction category api singleton instance
var (
	TransactionCategories = &TransactionCategoriesApi{
		categories: services.TransactionCategories,
	}
)

// CategoryListHandler returns transaction category list of current user
func (a *TransactionCategoriesApi) CategoryListHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryListReq models.TransactionCategoryListRequest
	err := c.ShouldBindQuery(&categoryListReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryListHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	categories, err := a.categories.GetAllCategoriesByUid(uid, categoryListReq.Type, categoryListReq.ParentId)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryListHandler] failed to get categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	return a.getTransactionCategoryListByTypeResponse(categories, categoryListReq.ParentId)
}

// CategoryGetHandler returns one specific transaction category of current user
func (a *TransactionCategoriesApi) CategoryGetHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryGetReq models.TransactionCategoryGetRequest
	err := c.ShouldBindQuery(&categoryGetReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	category, err := a.categories.GetCategoryByCategoryId(uid, categoryGetReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryGetHandler] failed to get category \"id:%d\" for user \"uid:%d\", because %s", categoryGetReq.Id, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	categoryResp := category.ToTransactionCategoryInfoResponse()

	return categoryResp, nil
}

// CategoryCreateHandler saves a new transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryCreateHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryCreateReq models.TransactionCategoryCreateRequest
	err := c.ShouldBindJSON(&categoryCreateReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	if categoryCreateReq.Type < models.CATEGORY_TYPE_INCOME || categoryCreateReq.Type > models.CATEGORY_TYPE_TRANSFER {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryCreateHandler] category type invalid, type is %d", categoryCreateReq.Type)
		return nil, errs.ErrTransactionCategoryTypeInvalid
	}

	uid := c.GetCurrentUid()

	if categoryCreateReq.ParentId > 0 {
		parentCategory, err := a.categories.GetCategoryByCategoryId(uid, categoryCreateReq.ParentId)

		if err != nil {
			log.ErrorfWithRequestId(c, "[transaction_categories.CategoryCreateHandler] failed to get parent category \"id:%d\" for user \"uid:%d\", because %s", categoryCreateReq.ParentId, uid, err.Error())
			return nil, errs.Or(err, errs.ErrOperationFailed)
		}

		if parentCategory == nil {
			log.WarnfWithRequestId(c, "[transaction_categories.CategoryCreateHandler] parent category \"id:%d\" does not exist for user \"uid:%d\"", categoryCreateReq.ParentId, uid)
			return nil, errs.ErrParentTransactionCategoryNotFound
		}

		if parentCategory.ParentCategoryId > 0 {
			log.WarnfWithRequestId(c, "[transaction_categories.CategoryCreateHandler] parent category \"id:%d\" has another parent category \"id:%d\" for user \"uid:%d\"", parentCategory.CategoryId, parentCategory.ParentCategoryId, uid)
			return nil, errs.ErrCannotAddToSecondaryTransactionCategory
		}
	}

	var maxOrderId int

	if categoryCreateReq.ParentId <= 0 {
		maxOrderId, err = a.categories.GetMaxDisplayOrder(uid, categoryCreateReq.Type)
	} else {
		maxOrderId, err = a.categories.GetMaxSubCategoryDisplayOrder(uid, categoryCreateReq.Type, categoryCreateReq.ParentId)
	}

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryCreateHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	category := a.createNewCategoryModel(uid, &categoryCreateReq, maxOrderId+1)

	err = a.categories.CreateCategory(category)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryCreateHandler] failed to create category \"id:%d\" for user \"uid:%d\", because %s", category.CategoryId, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_categories.CategoryCreateHandler] user \"uid:%d\" has created a new category \"id:%d\" successfully", uid, category.CategoryId)

	categoryResp := category.ToTransactionCategoryInfoResponse()

	return categoryResp, nil
}

// CategoryCreateBatchHandler saves some new transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryCreateBatchHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryCreateBatchReq models.TransactionCategoryCreateBatchRequest
	err := c.ShouldBindJSON(&categoryCreateBatchReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryCreateBatchHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	categoryTypeMaxOrderMap := make(map[models.TransactionCategoryType]int)
	categoriesMap := make(map[*models.TransactionCategory][]*models.TransactionCategory)
	categoriesMap[nil] = make([]*models.TransactionCategory, len(categoryCreateBatchReq.Categories))
	totalCount := 0

	for i := 0; i < len(categoryCreateBatchReq.Categories); i++ {
		categoryCreateReq := categoryCreateBatchReq.Categories[i]
		var maxOrderId, exists = categoryTypeMaxOrderMap[categoryCreateReq.Type]

		if !exists {
			maxOrderId, err = a.categories.GetMaxDisplayOrder(uid, categoryCreateReq.Type)

			if err != nil {
				log.ErrorfWithRequestId(c, "[transaction_categories.CategoryCreateBatchHandler] failed to get max display order for user \"uid:%d\", because %s", uid, err.Error())
				return nil, errs.ErrOperationFailed
			}
		}

		category := a.createNewCategoryModel(uid, &models.TransactionCategoryCreateRequest{
			Name:  categoryCreateReq.Name,
			Type:  categoryCreateReq.Type,
			Icon:  categoryCreateReq.Icon,
			Color: categoryCreateReq.Color,
		}, maxOrderId+1)

		categoriesMap[category] = make([]*models.TransactionCategory, len(categoryCreateReq.SubCategories))

		for j := 0; j < len(categoryCreateReq.SubCategories); j++ {
			subCategory := a.createNewCategoryModel(uid, categoryCreateReq.SubCategories[j], j+1)
			categoriesMap[category][j] = subCategory
			totalCount++
		}

		categoriesMap[nil][i] = category
		categoryTypeMaxOrderMap[categoryCreateReq.Type] = maxOrderId + 1
		totalCount++
	}

	categories, err := a.categories.CreateCategories(uid, categoriesMap)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryCreateBatchHandler] failed to create categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_categories.CategoryCreateBatchHandler] user \"uid:%d\" has created categoroies successfully", uid)

	return a.getTransactionCategoryListByTypeResponse(categories, 0)
}

// CategoryModifyHandler saves an existed transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryModifyHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryModifyReq models.TransactionCategoryModifyRequest
	err := c.ShouldBindJSON(&categoryModifyReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	category, err := a.categories.GetCategoryByCategoryId(uid, categoryModifyReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryModifyHandler] failed to get category \"id:%d\" for user \"uid:%d\", because %s", categoryModifyReq.Id, uid, err.Error())
		return nil, errs.ErrOperationFailed
	}

	newCategory := &models.TransactionCategory{
		CategoryId: category.CategoryId,
		Uid:        uid,
		Name:       categoryModifyReq.Name,
		Icon:       categoryModifyReq.Icon,
		Color:      categoryModifyReq.Color,
		Comment:    categoryModifyReq.Comment,
		Hidden:     categoryModifyReq.Hidden,
	}

	if newCategory.Name == category.Name &&
		newCategory.Icon == category.Icon &&
		newCategory.Color == category.Color &&
		newCategory.Comment == category.Comment &&
		newCategory.Hidden == category.Hidden {
		return nil, errs.ErrNothingWillBeUpdated
	}

	err = a.categories.ModifyCategory(newCategory)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryModifyHandler] failed to update category \"id:%d\" for user \"uid:%d\", because %s", categoryModifyReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_categories.CategoryModifyHandler] user \"uid:%d\" has updated category \"id:%d\" successfully", uid, categoryModifyReq.Id)

	newCategory.Type = category.Type
	newCategory.ParentCategoryId = category.ParentCategoryId
	newCategory.DisplayOrder = category.DisplayOrder
	categoryResp := newCategory.ToTransactionCategoryInfoResponse()

	return categoryResp, nil
}

// CategoryHideHandler hides an existed transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryHideHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryHideReq models.TransactionCategoryHideRequest
	err := c.ShouldBindJSON(&categoryHideReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryHideHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.categories.HideCategory(uid, []int64{categoryHideReq.Id}, categoryHideReq.Hidden)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryHideHandler] failed to hide category \"id:%d\" for user \"uid:%d\", because %s", categoryHideReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_categories.CategoryHideHandler] user \"uid:%d\" has hidden category \"id:%d\"", uid, categoryHideReq.Id)
	return true, nil
}

// CategoryMoveHandler moves display order of existed transaction categories by request parameters for current user
func (a *TransactionCategoriesApi) CategoryMoveHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryMoveReq models.TransactionCategoryMoveRequest
	err := c.ShouldBindJSON(&categoryMoveReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryMoveHandler] parse request failed, because %s", err.Error())
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

	err = a.categories.ModifyCategoryDisplayOrders(uid, categories)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryMoveHandler] failed to move categories for user \"uid:%d\", because %s", uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_categories.CategoryMoveHandler] user \"uid:%d\" has moved categories", uid)
	return true, nil
}

// CategoryDeleteHandler deletes an existed transaction category by request parameters for current user
func (a *TransactionCategoriesApi) CategoryDeleteHandler(c *core.Context) (interface{}, *errs.Error) {
	var categoryDeleteReq models.TransactionCategoryDeleteRequest
	err := c.ShouldBindJSON(&categoryDeleteReq)

	if err != nil {
		log.WarnfWithRequestId(c, "[transaction_categories.CategoryDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	err = a.categories.DeleteCategory(uid, categoryDeleteReq.Id)

	if err != nil {
		log.ErrorfWithRequestId(c, "[transaction_categories.CategoryDeleteHandler] failed to delete category \"id:%d\" for user \"uid:%d\", because %s", categoryDeleteReq.Id, uid, err.Error())
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	log.InfofWithRequestId(c, "[transaction_categories.CategoryDeleteHandler] user \"uid:%d\" has deleted category \"id:%d\"", uid, categoryDeleteReq.Id)
	return true, nil
}

func (a *TransactionCategoriesApi) createNewCategoryModel(uid int64, categoryCreateReq *models.TransactionCategoryCreateRequest, order int) *models.TransactionCategory {
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

		if categoryResp.ParentId <= models.LevelOneTransactionParentId {
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
		if parentId <= 0 && categoryResps[i].ParentId == models.LevelOneTransactionParentId {
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
