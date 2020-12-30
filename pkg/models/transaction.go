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

// TransactionDbType represents transaction type in database
type TransactionDbType byte

// Transaction db types
const (
	TRANSACTION_DB_TYPE_MODIFY_BALANCE TransactionDbType = 1
	TRANSACTION_DB_TYPE_INCOME         TransactionDbType = 2
	TRANSACTION_DB_TYPE_EXPENSE        TransactionDbType = 3
	TRANSACTION_DB_TYPE_TRANSFER_OUT   TransactionDbType = 4
	TRANSACTION_DB_TYPE_TRANSFER_IN    TransactionDbType = 5
)

// Transaction represents transaction data stored in database
type Transaction struct {
	TransactionId        int64             `xorm:"PK"`
	Uid                  int64             `xorm:"UNIQUE(UQE_transaction_uid_time) INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) NOT NULL"`
	Deleted              bool              `xorm:"INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) NOT NULL"`
	Type                 TransactionDbType `xorm:"INDEX(IDX_transaction_uid_deleted_type_time) NOT NULL"`
	CategoryId           int64             `xorm:"INDEX(IDX_transaction_uid_deleted_category_id_time) NOT NULL"`
	AccountId            int64             `xorm:"INDEX(IDX_transaction_uid_deleted_account_id_time) NOT NULL"`
	TransactionTime      int64             `xorm:"UNIQUE(UQE_transaction_uid_time) INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) NOT NULL"`
	Amount               int64             `xorm:"NOT NULL"`
	RelatedId            int64             `xorm:"NOT NULL"`
	RelatedAccountId     int64             `xorm:"NOT NULL"`
	RelatedAccountAmount int64             `xorm:"NOT NULL"`
	Comment              string            `xorm:"VARCHAR(255) NOT NULL"`
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
	DestinationAccountId int64           `json:"destinationAccountId,string" binding:"min=0"`
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
	DestinationAccountId int64    `json:"destinationAccountId,string" binding:"min=0"`
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
	DestinationAccountId int64           `json:"destinationAccountId,string,omitempty"`
	SourceAmount         int64           `json:"sourceAmount"`
	DestinationAmount    int64           `json:"destinationAmount,omitempty"`
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
	var transactionType TransactionType

	if c.Type == TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		transactionType = TRANSACTION_TYPE_MODIFY_BALANCE
	} else if c.Type == TRANSACTION_DB_TYPE_EXPENSE {
		transactionType = TRANSACTION_TYPE_EXPENSE
	} else if c.Type == TRANSACTION_DB_TYPE_INCOME {
		transactionType = TRANSACTION_TYPE_INCOME
	} else if c.Type == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		transactionType = TRANSACTION_TYPE_TRANSFER
	} else if c.Type == TRANSACTION_DB_TYPE_TRANSFER_IN {
		transactionType = TRANSACTION_TYPE_TRANSFER
	} else {
		return nil
	}

	sourceAccountId := c.AccountId
	sourceAmount := c.Amount

	destinationAccountId := int64(0)
	destinationAmount := int64(0)

	if c.Type == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		destinationAccountId = c.RelatedAccountId
		destinationAmount = c.RelatedAccountAmount
	} else if c.Type == TRANSACTION_DB_TYPE_TRANSFER_IN {
		sourceAccountId = c.RelatedAccountId
		sourceAmount = c.RelatedAccountAmount

		destinationAccountId = c.AccountId
		destinationAmount = c.Amount
	}

	return &TransactionInfoResponse{
		Id:                   c.TransactionId,
		TimeSequenceId:       c.TransactionTime,
		Type:                 transactionType,
		CategoryId:           c.CategoryId,
		Time:                 utils.GetUnixTimeFromTransactionTime(c.TransactionTime),
		SourceAccountId:      sourceAccountId,
		DestinationAccountId: destinationAccountId,
		SourceAmount:         sourceAmount,
		DestinationAmount:    destinationAmount,
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
