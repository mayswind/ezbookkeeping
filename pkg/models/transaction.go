package models

import "github.com/mayswind/lab/pkg/utils"

type TransactionType byte

const (
	TRANSACTION_TYPE_MODIFY_BALANCE TransactionType = 1
	TRANSACTION_TYPE_INCOME         TransactionType = 2
	TRANSACTION_TYPE_EXPENSE        TransactionType = 3
	TRANSACTION_TYPE_TRANSFER       TransactionType = 4
)

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

type TransactionCreateRequest struct {
	Type                 TransactionType `json:"type" binding:"required"`
	CategoryId           int64           `json:"categoryId,string"`
	Time                 int64           `json:"time" binding:"required,min=1"`
	SourceAccountId      int64           `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId int64           `json:"destinationAccountId,string" binding:"required,min=1"`
	SourceAmount         int64           `json:"sourceAmount"`
	DestinationAmount    int64           `json:"destinationAmount"`
	TagIds               []int64         `json:"tagIds,string"`
	Comment              string          `json:"comment" binding:"max=255"`
}

type TransactionModifyRequest struct {
	Id                   int64   `json:"id,string" binding:"required,min=1"`
	CategoryId           int64   `json:"categoryId,string"`
	Time                 int64   `json:"time" binding:"required,min=1"`
	SourceAccountId      int64   `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId int64   `json:"destinationAccountId,string" binding:"required,min=1"`
	SourceAmount         int64   `json:"sourceAmount"`
	DestinationAmount    int64   `json:"destinationAmount"`
	TagIds               []int64 `json:"tagIds,string"`
	Comment              string  `json:"comment" binding:"max=255"`
}

type TransactionListByMaxTimeRequest struct {
	MaxTime int64 `form:"max_time" binding:"required,min=1"`
	Count   int   `form:"count" binding:"required,min=1,max=50"`
}

type TransactionListInMonthByPageRequest struct {
	Year  int `form:"year" binding:"required,min=1"`
	Month int `form:"month" binding:"required,min=1"`
	Page  int `form:"page" binding:"required,min=1"`
	Count int `form:"count" binding:"required,min=1,max=50"`
}

type TransactionGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

type TransactionDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

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
	TagIds               []int64         `json:"tagIds,string"`
	Comment              string          `json:"comment"`
}

type TransactionInfoPageWrapperResponse struct {
	Items              TransactionInfoResponseSlice `json:"items"`
	NextTimeSequenceId *int64                       `json:"nextTimeSequenceId,string"`
}

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
		TagIds:               tagIds,
		Comment:              c.Comment,
	}
}

type TransactionInfoResponseSlice []*TransactionInfoResponse

func (c TransactionInfoResponseSlice) Len() int {
	return len(c)
}

func (c TransactionInfoResponseSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c TransactionInfoResponseSlice) Less(i, j int) bool {
	return c[i].Time < c[j].Time
}
