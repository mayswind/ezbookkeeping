package models

// ClearDataRequest represents all parameters of clear user data request
type ClearDataRequest struct {
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

// DataStatisticsResponse represents a view-object of user data statistic
type DataStatisticsResponse struct {
	TotalAccountCount              int64 `json:"totalAccountCount,string"`
	TotalTransactionCategoryCount  int64 `json:"totalTransactionCategoryCount,string"`
	TotalTransactionTagCount       int64 `json:"totalTransactionTagCount,string"`
	TotalTransactionCount          int64 `json:"totalTransactionCount,string"`
	TotalTransactionTemplateCount  int64 `json:"totalTransactionTemplateCount,string"`
	TotalScheduledTransactionCount int64 `json:"totalScheduledTransactionCount,string"`
}
