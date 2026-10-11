package models

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// Hours an archived (rejected) queued transaction is kept before being purged
const QueuedTransactionArchivedRetentionHours = 24

// QueuedTransactionStatus represents the status of a queued transaction in database
type QueuedTransactionStatus byte

// Queued transaction statuses
const (
	QUEUED_TRANSACTION_STATUS_PENDING  QueuedTransactionStatus = 1
	QUEUED_TRANSACTION_STATUS_ARCHIVED QueuedTransactionStatus = 2
)

// QueuedTransactionSource represents what produced a queued transaction
type QueuedTransactionSource byte

// Queued transaction sources
const (
	QUEUED_TRANSACTION_SOURCE_UNKNOWN            QueuedTransactionSource = 0
	QUEUED_TRANSACTION_SOURCE_SCHEDULED_TEMPLATE QueuedTransactionSource = 1
	QUEUED_TRANSACTION_SOURCE_API                QueuedTransactionSource = 2
	QUEUED_TRANSACTION_SOURCE_MCP                QueuedTransactionSource = 3
)

// QueuedTransaction represents a draft transaction waiting in the queue of the user.
//
// A queued transaction is not part of the bookkeeping data: it affects no account balance and shows
// up in no statistic. It only becomes real once the user accepts it, which creates a regular
// transaction out of it. Because it is a draft, every transaction field may be unset.
type QueuedTransaction struct {
	QueuedTransactionId  int64                   `xorm:"PK"`
	Uid                  int64                   `xorm:"INDEX(IDX_queued_transaction_uid_deleted_status_time) NOT NULL"`
	Deleted              bool                    `xorm:"INDEX(IDX_queued_transaction_uid_deleted_status_time) INDEX(IDX_queued_transaction_deleted_status_archived_time) NOT NULL"`
	Status               QueuedTransactionStatus `xorm:"INDEX(IDX_queued_transaction_uid_deleted_status_time) INDEX(IDX_queued_transaction_deleted_status_archived_time) NOT NULL"`
	Source               QueuedTransactionSource `xorm:"NOT NULL"`
	SourceTemplateId     int64                   `xorm:"NOT NULL"`
	Type                 TransactionType         `xorm:"NOT NULL"`
	CategoryId           int64                   `xorm:"NOT NULL"`
	AccountId            int64                   `xorm:"NOT NULL"`
	TransactionTime      int64                   `xorm:"INDEX(IDX_queued_transaction_uid_deleted_status_time) NOT NULL"`
	TimezoneUtcOffset    int16                   `xorm:"NOT NULL"`
	Amount               int64                   `xorm:"NOT NULL"`
	RelatedAccountId     int64                   `xorm:"NOT NULL"`
	RelatedAccountAmount int64                   `xorm:"NOT NULL"`
	HideAmount           bool                    `xorm:"NOT NULL"`
	TagIds               string                  `xorm:"VARCHAR(255) NOT NULL"`
	Comment              string                  `xorm:"VARCHAR(255) NOT NULL"`
	GeoLongitude         float64
	GeoLatitude          float64
	CreatedIp            string `xorm:"VARCHAR(39)"`
	CreatedUnixTime      int64
	UpdatedUnixTime      int64
	ArchivedUnixTime     int64 `xorm:"INDEX(IDX_queued_transaction_deleted_status_archived_time)"`
	DeletedUnixTime      int64
}

// QueuedTransactionListRequest represents all parameters of queued transaction list request
type QueuedTransactionListRequest struct {
	IncludeArchived bool `form:"includeArchived"`
}

// QueuedTransactionGetRequest represents all parameters of queued transaction getting request
type QueuedTransactionGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// QueuedTransactionCreateRequest represents all parameters of queued transaction creation request.
//
// Every field is optional on purpose: a draft is allowed to be incomplete, and the user fills in
// whatever is missing when accepting it.
type QueuedTransactionCreateRequest struct {
	Type                 TransactionType                `json:"type"`
	CategoryId           int64                          `json:"categoryId,string"`
	Time                 int64                          `json:"time"`
	UtcOffset            int16                          `json:"utcOffset" binding:"min=-720,max=840"`
	SourceAccountId      int64                          `json:"sourceAccountId,string"`
	DestinationAccountId int64                          `json:"destinationAccountId,string"`
	SourceAmount         int64                          `json:"sourceAmount" binding:"validTransactionAmount"`
	DestinationAmount    int64                          `json:"destinationAmount" binding:"validTransactionAmount"`
	HideAmount           bool                           `json:"hideAmount"`
	TagIds               []string                       `json:"tagIds"`
	Comment              string                         `json:"comment" binding:"max=255"`
	GeoLocation          *TransactionGeoLocationRequest `json:"geoLocation" binding:"omitempty"`
	ClientSessionId      string                         `json:"clientSessionId"`
}

// QueuedTransactionArchiveRequest represents all parameters of queued transaction rejection request
type QueuedTransactionArchiveRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// QueuedTransactionRestoreRequest represents all parameters of queued transaction restoration request
type QueuedTransactionRestoreRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// QueuedTransactionDeleteRequest represents all parameters of queued transaction deletion request
type QueuedTransactionDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// QueuedTransactionInfoResponse represents a view-object of a queued transaction
type QueuedTransactionInfoResponse struct {
	Id                   int64                            `json:"id,string"`
	Status               QueuedTransactionStatus          `json:"status"`
	Source               QueuedTransactionSource          `json:"source"`
	SourceTemplateId     int64                            `json:"sourceTemplateId,string,omitempty"`
	Type                 TransactionType                  `json:"type"`
	CategoryId           int64                            `json:"categoryId,string"`
	Category             *TransactionCategoryInfoResponse `json:"category,omitempty"`
	Time                 int64                            `json:"time"`
	UtcOffset            int16                            `json:"utcOffset"`
	SourceAccountId      int64                            `json:"sourceAccountId,string"`
	DestinationAccountId int64                            `json:"destinationAccountId,string,omitempty"`
	SourceAmount         int64                            `json:"sourceAmount"`
	DestinationAmount    *int64                           `json:"destinationAmount,omitempty"`
	HideAmount           bool                             `json:"hideAmount"`
	TagIds               []string                         `json:"tagIds"`
	Comment              string                           `json:"comment"`
	GeoLocation          *TransactionGeoLocationResponse  `json:"geoLocation,omitempty"`
	CreatedAt            int64                            `json:"createdAt"`
	ArchivedAt           int64                            `json:"archivedAt,omitempty"`
	ExpiresAt            int64                            `json:"expiresAt,omitempty"`
}

// QueuedTransactionCountResponse represents the response of the queued transaction count request
type QueuedTransactionCountResponse struct {
	PendingCount  int64 `json:"pendingCount"`
	ArchivedCount int64 `json:"archivedCount"`
}

// QueuedTransactionBatchOperationResponse represents the result of an operation over the whole queue
type QueuedTransactionBatchOperationResponse struct {
	AffectedCount int64 `json:"affectedCount"`
}

// GetTagIds returns all tag ids of the queued transaction
func (t *QueuedTransaction) GetTagIds() []int64 {
	tagIds := make([]string, 0)

	if t.TagIds != "" {
		tagIds = strings.Split(t.TagIds, ",")
	}

	result, _ := utils.StringArrayToInt64Array(tagIds)

	return result
}

// ToQueuedTransactionInfoResponse returns a view-object according to database model
func (t *QueuedTransaction) ToQueuedTransactionInfoResponse() *QueuedTransactionInfoResponse {
	tagIds := make([]string, 0, 0)

	if t.TagIds != "" {
		tagIds = strings.Split(t.TagIds, ",")
	}

	var destinationAmount *int64

	if t.Type == TRANSACTION_TYPE_TRANSFER {
		destinationAmount = &t.RelatedAccountAmount
	}

	var geoLocation *TransactionGeoLocationResponse

	if t.GeoLongitude != 0 || t.GeoLatitude != 0 {
		geoLocation = &TransactionGeoLocationResponse{
			Longitude: t.GeoLongitude,
			Latitude:  t.GeoLatitude,
		}
	}

	response := &QueuedTransactionInfoResponse{
		Id:                   t.QueuedTransactionId,
		Status:               t.Status,
		Source:               t.Source,
		SourceTemplateId:     t.SourceTemplateId,
		Type:                 t.Type,
		CategoryId:           t.CategoryId,
		Time:                 utils.GetUnixTimeFromTransactionTime(t.TransactionTime),
		UtcOffset:            t.TimezoneUtcOffset,
		SourceAccountId:      t.AccountId,
		DestinationAccountId: t.RelatedAccountId,
		SourceAmount:         t.Amount,
		DestinationAmount:    destinationAmount,
		HideAmount:           t.HideAmount,
		TagIds:               tagIds,
		Comment:              t.Comment,
		GeoLocation:          geoLocation,
		CreatedAt:            t.CreatedUnixTime,
	}

	if t.Status == QUEUED_TRANSACTION_STATUS_ARCHIVED && t.ArchivedUnixTime > 0 {
		response.ArchivedAt = t.ArchivedUnixTime
		response.ExpiresAt = t.ArchivedUnixTime + QueuedTransactionArchivedRetentionHours*60*60
	}

	return response
}

// QueuedTransactionInfoResponseSlice represents the slice data structure of QueuedTransactionInfoResponse
type QueuedTransactionInfoResponseSlice []*QueuedTransactionInfoResponse

// Len returns the count of items
func (s QueuedTransactionInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s QueuedTransactionInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s QueuedTransactionInfoResponseSlice) Less(i, j int) bool {
	if s[i].Status != s[j].Status {
		return s[i].Status < s[j].Status
	}

	if s[i].Time != s[j].Time {
		return s[i].Time > s[j].Time
	}

	return s[i].CreatedAt > s[j].CreatedAt
}
