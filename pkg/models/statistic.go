package models

// TransactionStatisticRequest represents all parameters of transaction statistic request
type TransactionStatisticRequest struct {
	StartTime int64 `form:"start_time" binding:"min=0"`
	EndTime   int64 `form:"end_time" binding:"min=0"`
}

// TransactionStatisticResponse represents an item of transaction overview
type TransactionStatisticResponse struct {
	StartTime int64                               `json:"startTime"`
	EndTime   int64                               `json:"endTime"`
	Items     []*TransactionStatisticResponseItem `json:"items"`
}

// TransactionStatisticResponseItem represents total amount item for an response
type TransactionStatisticResponseItem struct {
	CategoryId  int64 `json:"categoryId,string"`
	AccountId   int64 `json:"accountId,string"`
	TotalAmount int64 `json:"amount"`
}
