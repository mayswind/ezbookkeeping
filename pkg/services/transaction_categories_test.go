package services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func TestGetCategoryMapByList_EmptyList(t *testing.T) {
	categories := make([]*models.TransactionCategory, 0)
	actualCategoryMap := TransactionCategories.GetCategoryMapByList(categories)

	assert.NotNil(t, actualCategoryMap)
	assert.Equal(t, 0, len(actualCategoryMap))
}

func TestGetCategoryMapByList_MultipleCategories(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       1001,
			Name:             "Category Name",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
		{
			CategoryId:       1002,
			Name:             "Category Name2",
			Type:             models.CATEGORY_TYPE_INCOME,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
		{
			CategoryId:       1003,
			Name:             "Category Name3",
			Type:             models.CATEGORY_TYPE_TRANSFER,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           true,
		},
	}
	actualCategoryMap := TransactionCategories.GetCategoryMapByList(categories)

	assert.Equal(t, 3, len(actualCategoryMap))
	assert.Contains(t, actualCategoryMap, int64(1001))
	assert.Contains(t, actualCategoryMap, int64(1002))
	assert.Contains(t, actualCategoryMap, int64(1003))
	assert.Equal(t, "Category Name", actualCategoryMap[1001].Name)
	assert.Equal(t, "Category Name2", actualCategoryMap[1002].Name)
	assert.Equal(t, "Category Name3", actualCategoryMap[1003].Name)
}

func TestGetVisibleSubCategoryNameMapByList_EmptyList(t *testing.T) {
	categories := make([]*models.TransactionCategory, 0)
	expenseCategoryMap, incomeCategoryMap, transferCategoryMap := TransactionCategories.GetVisibleSubCategoryNameMapByList(categories)

	assert.NotNil(t, expenseCategoryMap)
	assert.NotNil(t, incomeCategoryMap)
	assert.NotNil(t, transferCategoryMap)
	assert.Equal(t, 0, len(expenseCategoryMap))
	assert.Equal(t, 0, len(incomeCategoryMap))
	assert.Equal(t, 0, len(transferCategoryMap))
}

func TestGetVisibleSubCategoryNameMapByList_OnlyParentCategories(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       1001,
			Name:             "Category Name",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
		{
			CategoryId:       1002,
			Name:             "Category Name2",
			Type:             models.CATEGORY_TYPE_INCOME,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
	}
	expenseCategoryMap, incomeCategoryMap, transferCategoryMap := TransactionCategories.GetVisibleSubCategoryNameMapByList(categories)

	assert.Equal(t, 0, len(expenseCategoryMap))
	assert.Equal(t, 0, len(incomeCategoryMap))
	assert.Equal(t, 0, len(transferCategoryMap))
}

func TestGetVisibleSubCategoryNameMapByList_WithHiddenCategories(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       1001,
			Name:             "Category Name",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
		{
			CategoryId:       2001,
			Name:             "Category Name2",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: 1001,
			Hidden:           true,
		},
		{
			CategoryId:       2002,
			Name:             "Category Name3",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: 1001,
			Hidden:           false,
		},
	}
	expenseCategoryMap, incomeCategoryMap, transferCategoryMap := TransactionCategories.GetVisibleSubCategoryNameMapByList(categories)

	assert.Equal(t, 1, len(expenseCategoryMap))
	assert.Contains(t, expenseCategoryMap, "Category Name3")
	assert.NotContains(t, expenseCategoryMap, "Category Name2")
	assert.Equal(t, 0, len(incomeCategoryMap))
	assert.Equal(t, 0, len(transferCategoryMap))
}

func TestGetVisibleSubCategoryNameMapByList_AllTypes(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       1001,
			Name:             "Category Name",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
		{
			CategoryId:       2001,
			Name:             "Category Name2",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: 1001,
			Hidden:           false,
		},
		{
			CategoryId:       1002,
			Name:             "Category Name3",
			Type:             models.CATEGORY_TYPE_INCOME,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
		{
			CategoryId:       2002,
			Name:             "Category Name4",
			Type:             models.CATEGORY_TYPE_INCOME,
			ParentCategoryId: 1002,
			Hidden:           false,
		},
		{
			CategoryId:       1003,
			Name:             "Category Name5",
			Type:             models.CATEGORY_TYPE_TRANSFER,
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
			Hidden:           false,
		},
		{
			CategoryId:       2003,
			Name:             "Category Name6",
			Type:             models.CATEGORY_TYPE_TRANSFER,
			ParentCategoryId: 1003,
			Hidden:           false,
		},
	}
	expenseCategoryMap, incomeCategoryMap, transferCategoryMap := TransactionCategories.GetVisibleSubCategoryNameMapByList(categories)

	assert.Equal(t, 1, len(expenseCategoryMap))
	assert.Contains(t, expenseCategoryMap, "Category Name2")
	assert.Contains(t, expenseCategoryMap["Category Name2"], "Category Name")

	assert.Equal(t, 1, len(incomeCategoryMap))
	assert.Contains(t, incomeCategoryMap, "Category Name4")
	assert.Contains(t, incomeCategoryMap["Category Name4"], "Category Name3")

	assert.Equal(t, 1, len(transferCategoryMap))
	assert.Contains(t, transferCategoryMap, "Category Name6")
	assert.Contains(t, transferCategoryMap["Category Name6"], "Category Name5")
}

func TestGetVisibleSubCategoryNameMapByList_OrphanSubCategories(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       2001,
			Name:             "Category Name",
			Type:             models.CATEGORY_TYPE_EXPENSE,
			ParentCategoryId: 9999,
			Hidden:           false,
		},
	}
	expenseCategoryMap, incomeCategoryMap, transferCategoryMap := TransactionCategories.GetVisibleSubCategoryNameMapByList(categories)

	assert.Equal(t, 0, len(expenseCategoryMap))
	assert.Equal(t, 0, len(incomeCategoryMap))
	assert.Equal(t, 0, len(transferCategoryMap))
}

func TestGetCategoryNames_EmptyList(t *testing.T) {
	categories := make([]*models.TransactionCategory, 0)
	actualNames := TransactionCategories.GetCategoryNames(categories)

	assert.NotNil(t, actualNames)
	assert.Equal(t, 0, len(actualNames))
}

func TestGetCategoryNames_MultipleCategories(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId: 1001,
			Name:       "Category Name",
		},
		{
			CategoryId: 1002,
			Name:       "Category Name2",
		},
		{
			CategoryId: 1003,
			Name:       "Category Name3",
		},
	}
	actualNames := TransactionCategories.GetCategoryNames(categories)

	assert.Equal(t, 3, len(actualNames))
	assert.Equal(t, "Category Name", actualNames[0])
	assert.Equal(t, "Category Name2", actualNames[1])
	assert.Equal(t, "Category Name3", actualNames[2])
}

func TestGetCategoryOrSubCategoryIdsByCategoryName_EmptyList(t *testing.T) {
	categories := make([]*models.TransactionCategory, 0)
	actualIds := TransactionCategories.GetCategoryOrSubCategoryIdsByCategoryName(categories, "Category Name")

	assert.NotNil(t, actualIds)
	assert.Equal(t, 0, len(actualIds))
}

func TestGetCategoryOrSubCategoryIdsByCategoryName_NotExistName(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       1001,
			Name:             "Category Name",
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
		},
	}
	actualIds := TransactionCategories.GetCategoryOrSubCategoryIdsByCategoryName(categories, "Non-existent Category")

	assert.NotNil(t, actualIds)
	assert.Equal(t, 0, len(actualIds))
}

func TestGetCategoryOrSubCategoryIdsByCategoryName_ParentCategoryWithoutChildren(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       1001,
			Name:             "Category Name",
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
		},
	}
	actualIds := TransactionCategories.GetCategoryOrSubCategoryIdsByCategoryName(categories, "Category Name")

	assert.NotNil(t, actualIds)
	assert.Equal(t, 0, len(actualIds))
}

func TestGetCategoryOrSubCategoryIdsByCategoryName_BothParentAndSubCategory(t *testing.T) {
	categories := []*models.TransactionCategory{
		{
			CategoryId:       1001,
			Name:             "Category Name",
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
		},
		{
			CategoryId:       2001,
			Name:             "Category Name",
			ParentCategoryId: 1001,
		},
		{
			CategoryId:       2002,
			Name:             "Category Name2",
			ParentCategoryId: 1001,
		},
		{
			CategoryId:       1002,
			Name:             "Category Name3",
			ParentCategoryId: models.LevelOneTransactionCategoryParentId,
		},
		{
			CategoryId:       2003,
			Name:             "Category Name",
			ParentCategoryId: 1002,
		},
	}
	actualIds := TransactionCategories.GetCategoryOrSubCategoryIdsByCategoryName(categories, "Category Name")

	assert.Equal(t, 3, len(actualIds))
	assert.Contains(t, actualIds, int64(2001))
	assert.Contains(t, actualIds, int64(2002))
	assert.Contains(t, actualIds, int64(2003))
}
