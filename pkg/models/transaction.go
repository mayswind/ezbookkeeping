package models

import "github.com/mayswind/lab/pkg/utils"

// TransactionType represents transaction type
type TransactionType byte

// Transaction types
const (
	TRANSACTION_TYPE_MODIFY_BALANCE TransactionType = 1
	TRANSACTION_TYPE_INCOME         TransactionType = 2
	TRANSACTION_TYPE_EXPENSE        TransactionType = 3
	TRANSACTION_TYPE_TRANSFER       TransactionType = 4
)

// Transaction represents transaction data stored in database
type Transaction struct {
	TransactionId        int64           `xorm:"PK"`
	Uid                  int64           `xorm:"UNIQUE(UQE_transaction_uid_transaction_time) INDEX(IDX_transaction_uid_deleted_transaction_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_category_id_time) NOT NULL"`
	Deleted              bool            `xorm:"INDEX(IDX_transaction_uid_deleted_transaction_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_category_id_time) NOT NULL"`
	Type                 TransactionType `xorm:"INDEX(IDX_transaction_uid_deleted_type_time) NOT NULL"`
	CategoryId           int64           `xorm:"INDEX(IDX_transaction_uid_deleted_category_id_time) NOT NULL"`
	TransactionTime      int64           `xorm:"UNIQUE(UQE_transaction_uid_transaction_time) INDEX(IDX_transaction_uid_deleted_transaction_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_category_id_time) NOT NULL"`
	SourceAccountId      int64           `xorm:"NOT NULL"`
	DestinationAccountId int64           `xorm:"NOT NULL"`
	SourceAmount         int64           `xorm:"NOT NULL"`
	DestinationAmount    int64           `xorm:"NOT NULL"`
	Comment              string          `xorm:"VARCHAR(255) NOT NULL"`
	CreatedUnixTime      int64
	UpdatedUnixTime      int64
	DeletedUnixTime      int64
}

// TransactionCreateRequest represents all parameters of transaction creation request
type TransactionCreateRequest struct {
	Type                 TransactionType `json:"type" binding:"required"`
	CategoryId           int64           `json:"categoryId,string"`
	Time                 int64           `json:"time" binding:"required,min=1"`
	SourceAccountId      int64           `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId int64           `json:"destinationAccountId,string" binding:"required,min=1"`
	SourceAmount         int64           `json:"sourceAmount" binding:"min=-99999999999,max=99999999999"`
	DestinationAmount    int64           `json:"destinationAmount" binding:"min=-99999999999,max=99999999999"`
	TagIds               []string        `json:"tagIds"`
	Comment              string          `json:"comment" binding:"max=255"`
}

// TransactionModifyRequest represents all parameters of transaction modification request
type TransactionModifyRequest struct {
	Id                   int64    `json:"id,string" binding:"required,min=1"`
	CategoryId           int64    `json:"categoryId,string"`
	Time                 int64    `json:"time" binding:"required,min=1"`
	SourceAccountId      int64    `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId int64    `json:"destinationAccountId,string" binding:"required,min=1"`
	SourceAmount         int64    `json:"sourceAmount" binding:"min=-99999999999,max=99999999999"`
	DestinationAmount    int64    `json:"destinationAmount" binding:"min=-99999999999,max=99999999999"`
	TagIds               []string `json:"tagIds"`
	Comment              string   `json:"comment" binding:"max=255"`
}

// TransactionListByMaxTimeRequest represents all parameters of transaction listing by max time request
type TransactionListByMaxTimeRequest struct {
	MaxTime int64 `form:"max_time" binding:"min=0"`
	Count   int   `form:"count" binding:"required,min=1,max=50"`
}

// TransactionListInMonthByPageRequest represents all parameters of transaction listing by month request
type TransactionListInMonthByPageRequest struct {
	Year  int `form:"year" binding:"required,min=1"`
	Month int `form:"month" binding:"required,min=1"`
	Page  int `form:"page" binding:"required,min=1"`
	Count int `form:"count" binding:"required,min=1,max=50"`
}

// TransactionGetRequest represents all parameters of transaction getting request
type TransactionGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// TransactionDeleteRequest represents all parameters of transaction deleting request
type TransactionDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// TransactionInfoResponse represents a view-object of transaction
type TransactionInfoResponse struct {
	Id                   int64           `json:"id,string"`
	TimeSequenceId       int64           `json:"timeSequenceId,string"`
	Type                 TransactionType `json:"type"`
	CategoryId           int64           `json:"categoryId,string"`
	Time                 int64           `json:"time"`
	SourceAccountId      int64           `json:"sourceAccountId,string"`
	DestinationAccountId int64           `json:"destinationAccountId,string"`
	SourceAmount         int64           `json:"sourceAmount"`
	DestinationAmount    int64           `json:"destinationAmount"`
	TagIds               []string        `json:"tagIds"`
	Comment              string          `json:"comment"`
}

// TransactionInfoPageWrapperResponse represents a response of transaction which contains items and next id
type TransactionInfoPageWrapperResponse struct {
	Items              TransactionInfoResponseSlice `json:"items"`
	NextTimeSequenceId *int64                       `json:"nextTimeSequenceId,string"`
}

// ToTransactionInfoResponse returns a view-object according to database model
func (c *Transaction) ToTransactionInfoResponse(tagIds []int64) *TransactionInfoResponse {
	return &TransactionInfoResponse{
		Id:                   c.TransactionId,
		TimeSequenceId:       c.TransactionTime,
		Type:                 c.Type,
		CategoryId:           c.CategoryId,
		Time:                 utils.GetUnixTimeFromTransactionTime(c.TransactionTime),
		SourceAccountId:      c.SourceAccountId,
		DestinationAccountId: c.DestinationAccountId,
		SourceAmount:         c.SourceAmount,
		DestinationAmount:    c.DestinationAmount,
		TagIds:               utils.Int64ArrayToStringArray(tagIds),
		Comment:              c.Comment,
	}
}

// TransactionInfoResponseSlice represents the slice data structure of TransactionInfoResponse
type TransactionInfoResponseSlice []*TransactionInfoResponse

// Len returns the count of items
func (c TransactionInfoResponseSlice) Len() int {
	return len(c)
}

// Swap swaps two items
func (c TransactionInfoResponseSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less reports whether the first item is less than the second one
func (c TransactionInfoResponseSlice) Less(i, j int) bool {
	if c[i].Time != c[j].Time {
		return c[i].Time > c[j].Time
	}

	return c[i].Id > c[j].Id
}
