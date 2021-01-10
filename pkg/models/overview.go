package models

// TransactionOverviewRequest represents all parameters of transaction overview request
type TransactionOverviewRequest struct {
	Query string `form:"query"`
}

// TransactionOverviewRequestItem represents an item of transaction overview request
type TransactionOverviewRequestItem struct {
	Name      string
	StartTime int64
	EndTime   int64
}

// TransactionOverviewResponseItem represents an item of transaction overview
type TransactionOverviewResponseItem struct {
	StartTime     int64 `json:"startTime"`
	EndTime       int64 `json:"endTime"`
	IncomeAmount  int64 `json:"incomeAmount"`
	ExpenseAmount int64 `json:"expenseAmount"`
}
