package models

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// TransactionTemplateType represents transaction template type in database
type TransactionTemplateType byte

// Transaction template types
const (
	TRANSACTION_TEMPLATE_TYPE_NORMAL   TransactionTemplateType = 1
	TRANSACTION_TEMPLATE_TYPE_SCHEDULE TransactionTemplateType = 2
)

// TransactionScheduleFrequencyType represents transaction template schedule frequency type
type TransactionScheduleFrequencyType byte

// Transaction template schedule frequency types
const (
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED TransactionScheduleFrequencyType = 0
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY   TransactionScheduleFrequencyType = 1
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY  TransactionScheduleFrequencyType = 2
)

// TransactionTemplate represents transaction template stored in database
type TransactionTemplate struct {
	TemplateId                 int64                            `xorm:"PK"`
	Uid                        int64                            `xorm:"INDEX(IDX_transaction_template_uid_deleted_template_type_order) NOT NULL"`
	Deleted                    bool                             `xorm:"INDEX(IDX_transaction_template_uid_deleted_template_type_order) INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_at) NOT NULL"`
	TemplateType               TransactionTemplateType          `xorm:"INDEX(IDX_transaction_template_uid_deleted_template_type_order) INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_at) NOT NULL"`
	Name                       string                           `xorm:"VARCHAR(32) NOT NULL"`
	Type                       TransactionType                  `xorm:"NOT NULL"`
	CategoryId                 int64                            `xorm:"NOT NULL"`
	AccountId                  int64                            `xorm:"NOT NULL"`
	ScheduledFrequencyType     TransactionScheduleFrequencyType `xorm:"INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_at)"`
	ScheduledFrequency         string                           `xorm:"VARCHAR(100)"`
	ScheduledAt                int16                            `xorm:"INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_at)"`
	ScheduledTimezoneUtcOffset int16
	TagIds                     string `xorm:"VARCHAR(255) NOT NULL"`
	Amount                     int64  `xorm:"NOT NULL"`
	RelatedAccountId           int64  `xorm:"NOT NULL"`
	RelatedAccountAmount       int64  `xorm:"NOT NULL"`
	HideAmount                 bool   `xorm:"NOT NULL"`
	Comment                    string `xorm:"VARCHAR(255) NOT NULL"`
	DisplayOrder               int32  `xorm:"INDEX(IDX_transaction_template_uid_deleted_template_type_order) NOT NULL"`
	Hidden                     bool   `xorm:"NOT NULL"`
	CreatedUnixTime            int64
	UpdatedUnixTime            int64
	DeletedUnixTime            int64
}

// TransactionTemplateListRequest represents all parameters of transaction template list request
type TransactionTemplateListRequest struct {
	TemplateType TransactionTemplateType `form:"templateType"`
}

// TransactionTemplateGetRequest represents all parameters of transaction template getting request
type TransactionTemplateGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// TransactionTemplateCreateRequest represents all parameters of transaction template creation request
type TransactionTemplateCreateRequest struct {
	TemplateType               TransactionTemplateType           `json:"templateType"`
	Name                       string                            `json:"name" binding:"required,notBlank,max=32"`
	Type                       TransactionType                   `json:"type" binding:"required"`
	CategoryId                 int64                             `json:"categoryId,string" binding:"required,min=1"`
	SourceAccountId            int64                             `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId       int64                             `json:"destinationAccountId,string" binding:"min=0"`
	SourceAmount               int64                             `json:"sourceAmount" binding:"min=-99999999999,max=99999999999"`
	DestinationAmount          int64                             `json:"destinationAmount" binding:"min=-99999999999,max=99999999999"`
	HideAmount                 bool                              `json:"hideAmount"`
	TagIds                     []string                          `json:"tagIds"`
	Comment                    string                            `json:"comment" binding:"max=255"`
	ScheduledFrequencyType     *TransactionScheduleFrequencyType `json:"scheduledFrequencyType" binding:"omitempty"`
	ScheduledFrequency         *string                           `json:"scheduledFrequency" binding:"omitempty"`
	ScheduledTimezoneUtcOffset *int16                            `json:"utcOffset" binding:"omitempty,min=-720,max=840"`
	ClientSessionId            string                            `json:"clientSessionId"`
}

// TransactionTemplateModifyNameRequest represents all parameters of transaction template name modification request
type TransactionTemplateModifyNameRequest struct {
	Id   int64  `json:"id,string" binding:"required,min=1"`
	Name string `json:"name" binding:"required,notBlank,max=32"`
}

// TransactionTemplateModifyRequest represents all parameters of transaction template modification request
type TransactionTemplateModifyRequest struct {
	Id                         int64                             `json:"id,string" binding:"required,min=1"`
	Name                       string                            `json:"name" binding:"required,notBlank,max=32"`
	Type                       TransactionType                   `json:"type" binding:"required"`
	CategoryId                 int64                             `json:"categoryId,string" binding:"required,min=1"`
	SourceAccountId            int64                             `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId       int64                             `json:"destinationAccountId,string" binding:"min=0"`
	SourceAmount               int64                             `json:"sourceAmount" binding:"min=-99999999999,max=99999999999"`
	DestinationAmount          int64                             `json:"destinationAmount" binding:"min=-99999999999,max=99999999999"`
	HideAmount                 bool                              `json:"hideAmount"`
	TagIds                     []string                          `json:"tagIds"`
	Comment                    string                            `json:"comment" binding:"max=255"`
	ScheduledFrequencyType     *TransactionScheduleFrequencyType `json:"scheduledFrequencyType" binding:"omitempty"`
	ScheduledFrequency         *string                           `json:"scheduledFrequency" binding:"omitempty"`
	ScheduledTimezoneUtcOffset *int16                            `json:"utcOffset" binding:"omitempty,min=-720,max=840"`
}

// TransactionTemplateHideRequest represents all parameters of transaction template hiding request
type TransactionTemplateHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

// TransactionTemplateMoveRequest represents all parameters of transaction template moving request
type TransactionTemplateMoveRequest struct {
	NewDisplayOrders []*TransactionTemplateNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// TransactionTemplateNewDisplayOrderRequest represents a data pair of id and display order
type TransactionTemplateNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// TransactionTemplateDeleteRequest represents all parameters of transaction template deleting request
type TransactionTemplateDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

type TransactionTemplateInfoResponse struct {
	*TransactionInfoResponse
	TemplateType           TransactionTemplateType           `json:"templateType"`
	Name                   string                            `json:"name"`
	ScheduledFrequencyType *TransactionScheduleFrequencyType `json:"scheduledFrequencyType,omitempty"`
	ScheduledFrequency     *string                           `json:"scheduledFrequency,omitempty"`
	ScheduledAt            *int16                            `json:"scheduledAt,omitempty"`
	DisplayOrder           int32                             `json:"displayOrder"`
	Hidden                 bool                              `json:"hidden"`
}

// GetTagIds returns all tag ids of the transaction template
func (t *TransactionTemplate) GetTagIds() []int64 {
	tagIds := make([]string, 0)

	if t.TagIds != "" {
		tagIds = strings.Split(t.TagIds, ",")
	}

	result, _ := utils.StringArrayToInt64Array(tagIds)

	return result
}

// ToTransactionTemplateInfoResponse returns a view-object according to database model
func (t *TransactionTemplate) ToTransactionTemplateInfoResponse(serverUtcOffset int16) *TransactionTemplateInfoResponse {
	utcOffset := serverUtcOffset

	if t.TemplateType == TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		utcOffset = t.ScheduledTimezoneUtcOffset
	}

	response := &TransactionTemplateInfoResponse{
		TransactionInfoResponse: t.toTransactionInfoResponse(utcOffset),
		TemplateType:            t.TemplateType,
		Name:                    t.Name,
		DisplayOrder:            t.DisplayOrder,
		Hidden:                  t.Hidden,
	}

	if t.TemplateType == TRANSACTION_TEMPLATE_TYPE_SCHEDULE {
		response.ScheduledFrequencyType = &t.ScheduledFrequencyType
		response.ScheduledFrequency = &t.ScheduledFrequency
		response.ScheduledAt = &t.ScheduledAt
	}

	return response
}

func (t *TransactionTemplate) toTransactionInfoResponse(utcOffset int16) *TransactionInfoResponse {
	tagIds := make([]string, 0, 0)

	if t.TagIds != "" {
		tagIds = strings.Split(t.TagIds, ",")
	}

	return &TransactionInfoResponse{
		Id:                   t.TemplateId,
		TimeSequenceId:       utils.GetMinTransactionTimeFromUnixTime(t.CreatedUnixTime),
		Type:                 t.Type,
		CategoryId:           t.CategoryId,
		Time:                 0,
		UtcOffset:            utcOffset,
		SourceAccountId:      t.AccountId,
		DestinationAccountId: t.RelatedAccountId,
		SourceAmount:         t.Amount,
		DestinationAmount:    t.RelatedAccountAmount,
		HideAmount:           t.HideAmount,
		TagIds:               tagIds,
		Comment:              t.Comment,
		GeoLocation:          nil,
		Editable:             true,
	}
}

// TransactionTemplateInfoResponseSlice represents the slice data structure of TransactionTemplateInfoResponse
type TransactionTemplateInfoResponseSlice []*TransactionTemplateInfoResponse

// Len returns the count of items
func (s TransactionTemplateInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s TransactionTemplateInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s TransactionTemplateInfoResponseSlice) Less(i, j int) bool {
	return s[i].DisplayOrder < s[j].DisplayOrder
}
