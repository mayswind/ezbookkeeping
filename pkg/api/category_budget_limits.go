package api

import (
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/duplicatechecker"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// [PLUGIN:budget] CategoryBudgetLimitsApi implements the category budget limit HTTP API.
// Self-contained handler module; see docs/PLUGIN_DESIGN.md.
type CategoryBudgetLimitsApi struct {
	ApiUsingConfig
	ApiUsingDuplicateChecker
	budgets    *services.CategoryBudgetLimitService
	categories *services.TransactionCategoryService
}

// Initialize a category budget limit api singleton instance
var (
	CategoryBudgetLimits = &CategoryBudgetLimitsApi{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		ApiUsingDuplicateChecker: ApiUsingDuplicateChecker{
			ApiUsingConfig: ApiUsingConfig{
				container: settings.Container,
			},
			container: duplicatechecker.Container,
		},
		budgets:    services.CategoryBudgetLimits,
		categories: services.TransactionCategories,
	}
)

// BudgetListByMonthHandler returns all category budget limits overlapping the given month
func (a *CategoryBudgetLimitsApi) BudgetListByMonthHandler(c *core.WebContext) (any, *errs.Error) {
	var listReq models.CategoryBudgetLimitListByMonthRequest

	if err := c.ShouldBindQuery(&listReq); err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetListByMonthHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetListByMonthHandler] cannot get client timezone, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	startUnix, endUnix := computeMonthBoundary(listReq.StartDate, clientTimezone)

	uid := c.GetCurrentUid()
	budgets, e := a.budgets.GetAllBudgetLimitsByMonth(c, uid, startUnix, endUnix)

	if e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetListByMonthHandler] failed to get budgets for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	result := make([]*models.CategoryBudgetLimitInfoResponse, 0, len(budgets))

	for i := 0; i < len(budgets); i++ {
		result = append(result, budgets[i].ToCategoryBudgetLimitInfoResponse())
	}

	return result, nil
}

// BudgetGetHandler returns one specific category budget limit
func (a *CategoryBudgetLimitsApi) BudgetGetHandler(c *core.WebContext) (any, *errs.Error) {
	var getReq models.CategoryBudgetLimitGetRequest

	if err := c.ShouldBindQuery(&getReq); err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetGetHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()
	budget, e := a.budgets.GetBudgetLimitByBudgetId(c, uid, getReq.Id)

	if e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetGetHandler] failed to get budget \"id:%d\" for user \"uid:%d\", because %s", getReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	return budget.ToCategoryBudgetLimitInfoResponse(), nil
}

// BudgetCreateHandler creates a new category budget limit for a month
func (a *CategoryBudgetLimitsApi) BudgetCreateHandler(c *core.WebContext) (any, *errs.Error) {
	var createReq models.CategoryBudgetLimitCreateRequest

	if err := c.ShouldBindJSON(&createReq); err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetCreateHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetCreateHandler] cannot get client timezone, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	startUnix, endUnix := computeMonthBoundary(createReq.StartDate, clientTimezone)

	uid := c.GetCurrentUid()

	// Verify the target category exists and belongs to the current user.
	category, e := a.categories.GetCategoryByCategoryId(c, uid, createReq.CategoryId)

	if e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetCreateHandler] failed to get category \"id:%d\" for user \"uid:%d\", because %s", createReq.CategoryId, uid, e.Error())
		return nil, errs.Or(e, errs.ErrCategoryBudgetLimitCategoryNotFound)
	}

	// Budgeting is meaningful for expense categories (mirrors Firefly III WITHDRAWAL-only rule).
	if category.Type != models.CATEGORY_TYPE_EXPENSE {
		log.Warnf(c, "[category_budget_limits.BudgetCreateHandler] category \"id:%d\" is not an expense category for user \"uid:%d\"", createReq.CategoryId, uid)
		return nil, errs.ErrCategoryBudgetLimitCategoryNotFound
	}

	budget := &models.CategoryBudgetLimit{
		Uid:        uid,
		CategoryId: createReq.CategoryId,
		StartDate:  startUnix,
		EndDate:    endUnix,
		Amount:     createReq.Amount,
		Currency:   createReq.Currency,
	}

	// Duplicate-submission guard (mirrors the category create handler).
	if a.CurrentConfig().EnableDuplicateSubmissionsCheck && createReq.ClientSessionId != "" {
		found, remark := a.GetSubmissionRemark(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CATEGORY, uid, createReq.ClientSessionId)

		if found {
			log.Infof(c, "[category_budget_limits.BudgetCreateHandler] another budget \"id:%s\" has been created for user \"uid:%d\"", remark, uid)

			if budgetId, parseErr := utils.StringToInt64(remark); parseErr == nil {
				if existing, getErr := a.budgets.GetBudgetLimitByBudgetId(c, uid, budgetId); getErr == nil {
					return existing.ToCategoryBudgetLimitInfoResponse(), nil
				}
			}
		}
	}

	if e := a.budgets.CreateBudgetLimit(c, budget); e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetCreateHandler] failed to create budget for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	log.Infof(c, "[category_budget_limits.BudgetCreateHandler] user \"uid:%d\" has created a new budget \"id:%d\" successfully", uid, budget.BudgetId)

	a.SetSubmissionRemarkIfEnable(duplicatechecker.DUPLICATE_CHECKER_TYPE_NEW_CATEGORY, uid, createReq.ClientSessionId, utils.Int64ToString(budget.BudgetId))

	return budget.ToCategoryBudgetLimitInfoResponse(), nil
}

// BudgetModifyHandler modifies an existed category budget limit (amount / currency only)
func (a *CategoryBudgetLimitsApi) BudgetModifyHandler(c *core.WebContext) (any, *errs.Error) {
	var modifyReq models.CategoryBudgetLimitModifyRequest

	if err := c.ShouldBindJSON(&modifyReq); err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetModifyHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	budget := &models.CategoryBudgetLimit{
		BudgetId: modifyReq.Id,
		Uid:      uid,
		Amount:   modifyReq.Amount,
		Currency: modifyReq.Currency,
	}

	if e := a.budgets.ModifyBudgetLimit(c, budget); e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetModifyHandler] failed to modify budget \"id:%d\" for user \"uid:%d\", because %s", modifyReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	updated, e := a.budgets.GetBudgetLimitByBudgetId(c, uid, modifyReq.Id)

	if e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetModifyHandler] failed to reload budget \"id:%d\" for user \"uid:%d\", because %s", modifyReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	log.Infof(c, "[category_budget_limits.BudgetModifyHandler] user \"uid:%d\" has modified budget \"id:%d\" successfully", uid, modifyReq.Id)

	return updated.ToCategoryBudgetLimitInfoResponse(), nil
}

// BudgetDeleteHandler soft-deletes a category budget limit
func (a *CategoryBudgetLimitsApi) BudgetDeleteHandler(c *core.WebContext) (any, *errs.Error) {
	var deleteReq models.CategoryBudgetLimitDeleteRequest

	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetDeleteHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	uid := c.GetCurrentUid()

	if e := a.budgets.DeleteBudgetLimit(c, uid, deleteReq.Id); e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetDeleteHandler] failed to delete budget \"id:%d\" for user \"uid:%d\", because %s", deleteReq.Id, uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	log.Infof(c, "[category_budget_limits.BudgetDeleteHandler] user \"uid:%d\" has deleted budget \"id:%d\" successfully", uid, deleteReq.Id)

	return true, nil
}

// BudgetOverviewHandler returns the budget overview for a given month: each limit plus the
// actual amount spent on that category in the same period, plus totals.
func (a *CategoryBudgetLimitsApi) BudgetOverviewHandler(c *core.WebContext) (any, *errs.Error) {
	var listReq models.CategoryBudgetLimitListByMonthRequest

	if err := c.ShouldBindQuery(&listReq); err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetOverviewHandler] parse request failed, because %s", err.Error())
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	clientTimezone, err := c.GetClientTimezone()

	if err != nil {
		log.Warnf(c, "[category_budget_limits.BudgetOverviewHandler] cannot get client timezone, because %s", err.Error())
		return nil, errs.ErrClientTimezoneOffsetInvalid
	}

	startUnix, endUnix := computeMonthBoundary(listReq.StartDate, clientTimezone)

	uid := c.GetCurrentUid()
	overview, e := a.budgets.GetBudgetOverview(c, uid, startUnix, endUnix, clientTimezone, false)

	if e != nil {
		log.Errorf(c, "[category_budget_limits.BudgetOverviewHandler] failed to get overview for user \"uid:%d\", because %s", uid, e.Error())
		return nil, errs.Or(e, errs.ErrOperationFailed)
	}

	return overview, nil
}

// computeMonthBoundary returns the unix timestamps for 00:00:00 of the 1st and 23:59:59 of
// the last day of the month that contains the given anyMoment, interpreted in the client
// timezone. This guarantees two budgets for the same month produce an identical start_date
// key (used by the unique constraint).
func computeMonthBoundary(anyMomentUnix int64, clientTimezone *time.Location) (startUnix int64, endUnix int64) {
	moment := time.Unix(anyMomentUnix, 0).In(clientTimezone)
	firstOfMonth := time.Date(moment.Year(), moment.Month(), 1, 0, 0, 0, 0, clientTimezone)
	startUnix = firstOfMonth.Unix()
	nextMonth := firstOfMonth.AddDate(0, 1, 0)
	endUnix = nextMonth.Unix() - 1
	return
}
