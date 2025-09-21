package models

// ClearDataRequest represents all parameters of clear user data request
type ClearDataRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

// ClearAccountTransactionsRequest represents all parameters of clear transaction data of a specific account request
type ClearAccountTransactionsRequest struct {
	AccountId int64  `json:"accountId,string" binding:"required,min=1"`
	Password  string `json:"password" binding:"omitempty,min=6,max=128"`
}

// DataStatisticsResponse represents a view-object of user data statistic
type DataStatisticsResponse struct {
	TotalAccountCount              int64 `json:"totalAccountCount,string"`
	TotalTransactionCategoryCount  int64 `json:"totalTransactionCategoryCount,string"`
	TotalTransactionTagCount       int64 `json:"totalTransactionTagCount,string"`
	TotalTransactionCount          int64 `json:"totalTransactionCount,string"`
	TotalTransactionPictureCount   int64 `json:"totalTransactionPictureCount,string"`
	TotalTransactionTemplateCount  int64 `json:"totalTransactionTemplateCount,string"`
	TotalScheduledTransactionCount int64 `json:"totalScheduledTransactionCount,string"`
}

// ExportTransactionDataRequest represents export transaction request
type ExportTransactionDataRequest struct {
	Type          TransactionType          `form:"type" binding:"min=0,max=4"`
	CategoryIds   string                   `form:"category_ids"`
	AccountIds    string                   `form:"account_ids"`
	TagIds        string                   `form:"tag_ids"`
	TagFilterType TransactionTagFilterType `form:"tag_filter_type" binding:"min=0,max=3"`
	AmountFilter  string                   `form:"amount_filter" binding:"validAmountFilter"`
	Keyword       string                   `form:"keyword"`
	MaxTime       int64                    `form:"max_time" binding:"min=0"` // Unix timestamp in seconds
	MinTime       int64                    `form:"min_time" binding:"min=0"` // Unix timestamp in seconds
}
