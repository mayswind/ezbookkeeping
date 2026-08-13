package services

import (
	"math/big"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

// [PLUGIN:budget] CategoryBudgetLimitService implements category-based monthly budgeting.
// It is a self-contained service: it does NOT modify the core Transaction model nor
// TransactionService. The "actual spent" computation reuses the existing aggregation
// method services.Transactions.GetAccountsAndCategoriesTotalInflowAndOutflow, then filters
// in-memory by category and expense type. See docs/PLUGIN_DESIGN.md.
type CategoryBudgetLimitService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

// Initialize a category budget limit service singleton instance
var (
	CategoryBudgetLimits = &CategoryBudgetLimitService{
		ServiceUsingDB: ServiceUsingDB{
			container: datastore.Container,
		},
		ServiceUsingUuid: ServiceUsingUuid{
			container: uuid.Container,
		},
	}
)

// GetBudgetLimitByBudgetId returns a single category budget limit by its primary key
func (s *CategoryBudgetLimitService) GetBudgetLimitByBudgetId(c core.Context, uid int64, budgetId int64) (*models.CategoryBudgetLimit, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if budgetId <= 0 {
		return nil, errs.ErrCategoryBudgetLimitIdInvalid
	}

	budget := &models.CategoryBudgetLimit{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(budgetId).Where("uid=? AND deleted=?", uid, false).Get(budget)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrCategoryBudgetLimitNotFound
	}

	return budget, nil
}

// GetAllBudgetLimitsByMonth returns all category budget limits whose period overlaps the
// provided [startUnix, endUnix] window. Typically startUnix/endUnix is one calendar month.
func (s *CategoryBudgetLimitService) GetAllBudgetLimitsByMonth(c core.Context, uid int64, startUnix int64, endUnix int64) ([]*models.CategoryBudgetLimit, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	var budgets []*models.CategoryBudgetLimit
	err := s.UserDataDB(uid).NewSession(c).
		Where("uid=? AND deleted=? AND start_date<=? AND end_date>=?", uid, false, endUnix, startUnix).
		OrderBy("start_date asc, category_id asc").
		Find(&budgets)

	return budgets, err
}

// CreateBudgetLimit saves a new category budget limit. The start date / end date are
// normalized to the first second / last second of the month in the *client* timezone,
// so that two limits created for "the same month" produce an identical (uid, category_id,
// start_date) key and the unique constraint enforces one-limit-per-category-per-month.
func (s *CategoryBudgetLimitService) CreateBudgetLimit(c core.Context, budget *models.CategoryBudgetLimit) error {
	if budget.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if budget.CategoryId <= 0 {
		return errs.ErrCategoryBudgetLimitCategoryIdEmpty
	}

	if budget.Amount <= 0 {
		return errs.ErrCategoryBudgetLimitAmountInvalid
	}

	budget.BudgetId = s.GenerateUuid(uuid.UUID_TYPE_BUDGET)

	if budget.BudgetId < 1 {
		return errs.ErrSystemIsBusy
	}

	budget.Deleted = false
	budget.CreatedUnixTime = time.Now().Unix()
	budget.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(budget.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(budget)
		return err
	})
}

// ModifyBudgetLimit saves an existed category budget limit (only the mutable fields:
// amount and currency). The category and period are immutable after creation.
func (s *CategoryBudgetLimitService) ModifyBudgetLimit(c core.Context, budget *models.CategoryBudgetLimit) error {
	if budget.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if budget.Amount <= 0 {
		return errs.ErrCategoryBudgetLimitAmountInvalid
	}

	budget.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(budget.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(budget.BudgetId).Cols("amount", "currency", "updated_unix_time").
			Where("uid=? AND deleted=?", budget.Uid, false).Update(budget)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrCategoryBudgetLimitNotFound
		}

		return nil
	})
}

// DeleteBudgetLimit soft-deletes a category budget limit
func (s *CategoryBudgetLimitService) DeleteBudgetLimit(c core.Context, uid int64, budgetId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	if budgetId <= 0 {
		return errs.ErrCategoryBudgetLimitIdInvalid
	}

	now := time.Now().Unix()
	updateModel := &models.CategoryBudgetLimit{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		updatedRows, err := sess.ID(budgetId).Cols("deleted", "deleted_unix_time").
			Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrCategoryBudgetLimitNotFound
		}

		return nil
	})
}

// GetBudgetOverview returns the full budget overview for a month: for each limit in the
// period it computes how much has actually been spent on that category (including its
// sub-categories) in the same period, using the existing core aggregation method.
//
// clientTimezone is used both to scope the aggregation date window and to interpret the
// stored start/end unix times. useTransactionTimezone mirrors the statistics pages: when
// true, each transaction's own timezone is used for bucketing.
func (s *CategoryBudgetLimitService) GetBudgetOverview(c core.Context, uid int64, startUnix int64, endUnix int64, clientTimezone *time.Location, useTransactionTimezone bool) (*models.CategoryBudgetOverviewResponse, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	// 1. Load all limits in this month.
	budgets, err := s.GetAllBudgetLimitsByMonth(c, uid, startUnix, endUnix)

	if err != nil {
		return nil, err
	}

	if len(budgets) == 0 {
		return &models.CategoryBudgetOverviewResponse{
			StartDate: startUnix,
			EndDate:   endUnix,
		}, nil
	}

	// 2. Load all category names for display + expand parent categories to sub-category ids.
	allCategoryIds := make([]int64, 0, len(budgets))
	for i := 0; i < len(budgets); i++ {
		allCategoryIds = append(allCategoryIds, budgets[i].CategoryId)
	}

	categoryMap, err := TransactionCategories.GetCategoriesByCategoryIds(c, uid, allCategoryIds)

	if err != nil {
		return nil, err
	}

	// 3. Reuse the existing core aggregation: returns []*TransactionTotalAmount, one entry
	//    per (category_id, account_id) combination, with Amount as a *big.Int in minor units.
	totalAmounts, err := Transactions.GetAccountsAndCategoriesTotalInflowAndOutflow(c, uid, startUnix, endUnix, nil, false, "", core.MATCH_MODE_IGNORE_CASE, clientTimezone, useTransactionTimezone)

	if err != nil {
		return nil, err
	}

	// 4. Build a per-category expense total. Only EXPENSE rows count as "spent" (mirrors
	//    Firefly III's WITHDRAWAL-only budgeting). Transfer rows are ignored.
	categoryExpenseTotals := make(map[int64]*big.Int)

	for i := 0; i < len(totalAmounts); i++ {
		ta := totalAmounts[i]

		if ta.Type != models.TRANSACTION_DB_TYPE_EXPENSE {
			continue
		}

		if ta.CategoryId <= 0 {
			continue
		}

		if existing, ok := categoryExpenseTotals[ta.CategoryId]; ok {
			existing.Add(existing, ta.Amount)
		} else {
			categoryExpenseTotals[ta.CategoryId] = new(big.Int).Set(ta.Amount)
		}
	}

	// 5. Assemble the overview. For each budget, sum the spent amount across the budget
	//    category and (if the budget is on a parent category) its sub-categories.
	items := make([]*models.CategoryBudgetOverviewItem, 0, len(budgets))
	var totalLimit, totalActual int64

	for i := 0; i < len(budgets); i++ {
		budget := budgets[i]
		category, categoryExists := categoryMap[budget.CategoryId]

		// Determine which category ids contribute to this budget's spent total.
		contributingCategoryIds := []int64{budget.CategoryId}

		// If the budget is set on a parent (level-one) category, expand to all sub-categories.
		// We can detect a parent category by ParentCategoryId == 0, but to be robust also
		// expand when the category itself is a parent that has children: this requires
		// querying sub-categories. To avoid an N+1 query we only expand when the stored
		// category is a parent (ParentCategoryId == 0). Sub-category budgets simply use
		// their own id.
		if categoryExists && category.ParentCategoryId == models.LevelOneTransactionCategoryParentId {
			subCategories, subErr := TransactionCategories.GetSubCategoriesByCategoryIds(c, uid, []int64{budget.CategoryId})

			if subErr == nil {
				for j := 0; j < len(subCategories); j++ {
					contributingCategoryIds = append(contributingCategoryIds, subCategories[j].CategoryId)
				}
			}
		}

		spent := new(big.Int)
		for j := 0; j < len(contributingCategoryIds); j++ {
			if t, ok := categoryExpenseTotals[contributingCategoryIds[j]]; ok {
				spent.Add(spent, t)
			}
		}

		spentInt64 := spent.Int64()

		item := &models.CategoryBudgetOverviewItem{
			CategoryBudgetLimitInfoResponse: budget.ToCategoryBudgetLimitInfoResponse(),
			ActualExpenseAmount:             spentInt64,
			AvailableAmount:                 budget.Amount - spentInt64,
		}

		if categoryExists {
			item.CategoryName = category.Name
			item.CategoryParentId = category.ParentCategoryId
			item.CategoryType = category.Type
		}

		items = append(items, item)
		totalLimit += budget.Amount
		totalActual += spentInt64
	}

	return &models.CategoryBudgetOverviewResponse{
		StartDate:      startUnix,
		EndDate:        endUnix,
		TotalLimit:     totalLimit,
		TotalActual:    totalActual,
		TotalAvailable: totalLimit - totalActual,
		Items:          items,
	}, nil
}
