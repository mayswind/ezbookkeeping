package api

import (
	"testing"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestTransactionCategoryBudgetValidation(t *testing.T) {
	api := &TransactionCategoriesApi{}

	assert.True(t, api.isTransactionCategoryBudgetValid(models.CATEGORY_TYPE_EXPENSE, 0, ""))
	assert.True(t, api.isTransactionCategoryBudgetValid(models.CATEGORY_TYPE_EXPENSE, 10000, "RUB"))

	assert.False(t, api.isTransactionCategoryBudgetValid(models.CATEGORY_TYPE_EXPENSE, -1, "RUB"))
	assert.False(t, api.isTransactionCategoryBudgetValid(models.CATEGORY_TYPE_EXPENSE, 0, "RUB"))
	assert.False(t, api.isTransactionCategoryBudgetValid(models.CATEGORY_TYPE_EXPENSE, 10000, ""))
	assert.False(t, api.isTransactionCategoryBudgetValid(models.CATEGORY_TYPE_INCOME, 10000, "RUB"))
	assert.False(t, api.isTransactionCategoryBudgetValid(models.CATEGORY_TYPE_TRANSFER, 10000, "RUB"))
}
