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

// TransactionTotalAmount represents total amount for specific transaction type
type TransactionTotalAmount struct {
	Uid         int64
	Type        TransactionDbType
	TotalAmount int64 `xorm:"NOT NULL"`
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
	Type       TransactionDbType `form:"type" binding:"min=0,max=4"`
	CategoryId int64             `form:"category_id" binding:"min=0"`
	AccountId  int64             `form:"account_id" binding:"min=0"`
	Keyword    string            `form:"keyword"`
	MaxTime    int64             `form:"max_time" binding:"min=0"`
	MinTime    int64             `form:"min_time" binding:"min=0"`
	Count      int               `form:"count" binding:"required,min=1,max=50"`
}

// TransactionListInMonthByPageRequest represents all parameters of transaction listing by month request
type TransactionListInMonthByPageRequest struct {
	Year       int               `form:"year" binding:"required,min=1"`
	Month      int               `form:"month" binding:"required,min=1"`
	Type       TransactionDbType `form:"type" binding:"min=0,max=4"`
	CategoryId int64             `form:"category_id" binding:"min=0"`
	AccountId  int64             `form:"account_id" binding:"min=0"`
	Keyword    string            `form:"keyword"`
	Page       int               `form:"page" binding:"required,min=1"`
	Count      int               `form:"count" binding:"required,min=1,max=50"`
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
	Editable             bool            `json:"editable"`
}

// TransactionInfoPageWrapperResponse represents a response of transaction which contains items and next id
type TransactionInfoPageWrapperResponse struct {
	Items              TransactionInfoResponseSlice `json:"items"`
	NextTimeSequenceId *int64                       `json:"nextTimeSequenceId,string"`
}

// IsEditable returns whether this transaction can be edited
func (t *Transaction) IsEditable(account *Account, relatedAccount *Account) bool {
	if account == nil || account.Hidden {
		return false
	}

	if t.Type == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		if relatedAccount == nil || relatedAccount.Hidden {
			return false
		}
	}

	return true
}

// ToTransactionInfoResponse returns a view-object according to database model
func (t *Transaction) ToTransactionInfoResponse(tagIds []int64, editable bool) *TransactionInfoResponse {
	var transactionType TransactionType

	if t.Type == TRANSACTION_DB_TYPE_MODIFY_BALANCE {
		transactionType = TRANSACTION_TYPE_MODIFY_BALANCE
	} else if t.Type == TRANSACTION_DB_TYPE_EXPENSE {
		transactionType = TRANSACTION_TYPE_EXPENSE
	} else if t.Type == TRANSACTION_DB_TYPE_INCOME {
		transactionType = TRANSACTION_TYPE_INCOME
	} else if t.Type == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		transactionType = TRANSACTION_TYPE_TRANSFER
	} else if t.Type == TRANSACTION_DB_TYPE_TRANSFER_IN {
		transactionType = TRANSACTION_TYPE_TRANSFER
	} else {
		return nil
	}

	sourceAccountId := t.AccountId
	sourceAmount := t.Amount

	destinationAccountId := int64(0)
	destinationAmount := int64(0)

	if t.Type == TRANSACTION_DB_TYPE_TRANSFER_OUT {
		destinationAccountId = t.RelatedAccountId
		destinationAmount = t.RelatedAccountAmount
	} else if t.Type == TRANSACTION_DB_TYPE_TRANSFER_IN {
		sourceAccountId = t.RelatedAccountId
		sourceAmount = t.RelatedAccountAmount

		destinationAccountId = t.AccountId
		destinationAmount = t.Amount
	}

	return &TransactionInfoResponse{
		Id:                   t.TransactionId,
		TimeSequenceId:       t.TransactionTime,
		Type:                 transactionType,
		CategoryId:           t.CategoryId,
		Time:                 utils.GetUnixTimeFromTransactionTime(t.TransactionTime),
		SourceAccountId:      sourceAccountId,
		DestinationAccountId: destinationAccountId,
		SourceAmount:         sourceAmount,
		DestinationAmount:    destinationAmount,
		TagIds:               utils.Int64ArrayToStringArray(tagIds),
		Comment:              t.Comment,
		Editable:             editable,
	}
}

// TransactionInfoResponseSlice represents the slice data structure of TransactionInfoResponse
type TransactionInfoResponseSlice []*TransactionInfoResponse

// Len returns the count of items
func (s TransactionInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionInfoResponseSlice) Less(i, j int) bool {
	if s[i].Time != s[j].Time {
		return s[i].Time > s[j].Time
	}

	return s[i].Id > s[j].Id
}
