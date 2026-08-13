package models

// [PLUGIN:budget] This file implements the category budget limit feature as a self-contained
// extension. It does NOT modify the core Transaction model. See docs/PLUGIN_DESIGN.md.

// CategoryBudgetLimit represents a monthly spending limit set on a transaction category.
// One row per (uid, category_id, start_date). Mirrors Firefly III's budget_limits concept,
// simplified: no rollover, no join table (the link to transactions is implicit through
// Transaction.CategoryId, which already exists in the core model).
type CategoryBudgetLimit struct {
	BudgetId        int64  `xorm:"PK"`
	Uid             int64  `xorm:"INDEX(IDX_category_budget_uid_deleted) NOT NULL unique(UQE_category_budget_uid_category_startdate)"`
	CategoryId      int64  `xorm:"NOT NULL unique(UQE_category_budget_uid_category_startdate)"`
	StartDate       int64  `xorm:"NOT NULL unique(UQE_category_budget_uid_category_startdate)"` // unix time at 00:00:00 of the 1st of the month (client timezone)
	EndDate         int64  `xorm:"NOT NULL"`                                                  // unix time at 23:59:59 of the last day of the month (client timezone)
	Amount          int64  `xorm:"NOT NULL"`                                                  // limit in minor units (cents)
	Currency        string `xorm:"VARCHAR(8) NOT NULL"`                                       // base currency code for v1
	Deleted         bool   `xorm:"INDEX(IDX_category_budget_uid_deleted) NOT NULL"`
	CreatedUnixTime int64
	UpdatedUnixTime int64
	DeletedUnixTime int64
}

// CategoryBudgetLimitListByMonthRequest represents all parameters of category budget limit listing request
type CategoryBudgetLimitListByMonthRequest struct {
	// StartDate is the unix time at 00:00:00 of the 1st of the target month (client timezone)
	StartDate int64 `form:"startDate,string" binding:"required,min=1"`
}

// CategoryBudgetLimitGetRequest represents all parameters of a single category budget limit getting request
type CategoryBudgetLimitGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// CategoryBudgetLimitCreateRequest represents all parameters of a single category budget limit creation request
type CategoryBudgetLimitCreateRequest struct {
	CategoryId      int64  `json:"categoryId,string" binding:"required,min=1"`
	StartDate       int64  `json:"startDate,string" binding:"required,min=1"` // unix time at 00:00:00 of the 1st of the month
	Amount          int64  `json:"amount" binding:"required,validTransactionAmount"`
	Currency        string `json:"currency" binding:"required,max=8"`
	ClientSessionId string `json:"clientSessionId"`
}

// CategoryBudgetLimitModifyRequest represents all parameters of a single category budget limit modification request
type CategoryBudgetLimitModifyRequest struct {
	Id         int64 `json:"id,string" binding:"required,min=1"`
	Amount     int64 `json:"amount" binding:"required,validTransactionAmount"`
	Currency   string `json:"currency" binding:"required,max=8"`
}

// CategoryBudgetLimitDeleteRequest represents all parameters of a single category budget limit deletion request
type CategoryBudgetLimitDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// CategoryBudgetLimitInfoResponse represents a view-object of a single category budget limit
type CategoryBudgetLimitInfoResponse struct {
	Id              int64  `json:"id,string"`
	CategoryId      int64  `json:"categoryId,string"`
	StartDate       int64  `json:"startDate,string"`
	EndDate         int64  `json:"endDate,string"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
}

// ToCategoryBudgetLimitInfoResponse returns a view-object according to the database model
func (b *CategoryBudgetLimit) ToCategoryBudgetLimitInfoResponse() *CategoryBudgetLimitInfoResponse {
	return &CategoryBudgetLimitInfoResponse{
		Id:         b.BudgetId,
		CategoryId: b.CategoryId,
		StartDate:  b.StartDate,
		EndDate:    b.EndDate,
		Amount:     b.Amount,
		Currency:   b.Currency,
	}
}

// CategoryBudgetOverviewItem represents one row of the budget overview for a given month:
// the limit set on a category plus how much has actually been spent on that category in the period.
type CategoryBudgetOverviewItem struct {
	*CategoryBudgetLimitInfoResponse
	CategoryName      string `json:"categoryName"`
	CategoryParentId  int64  `json:"categoryParentId,string"`
	CategoryType      TransactionCategoryType `json:"categoryType"`
	ActualExpenseAmount int64 `json:"actualExpenseAmount"` // total spent in the period, minor units (already base currency)
	AvailableAmount   int64 `json:"availableAmount"`       // Amount - ActualExpenseAmount (can be negative)
}

// CategoryBudgetOverviewResponse represents the budget overview for a given month
type CategoryBudgetOverviewResponse struct {
	StartDate       int64                        `json:"startDate,string"`
	EndDate         int64                        `json:"endDate,string"`
	TotalLimit      int64                        `json:"totalLimit"`
	TotalActual     int64                        `json:"totalActualExpenseAmount"`
	TotalAvailable  int64                        `json:"totalAvailableAmount"`
	Items           []*CategoryBudgetOverviewItem `json:"items"`
}
