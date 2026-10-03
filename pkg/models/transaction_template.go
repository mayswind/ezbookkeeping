package models

import (
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
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
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED     TransactionScheduleFrequencyType = 0
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY       TransactionScheduleFrequencyType = 1
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY      TransactionScheduleFrequencyType = 2
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY        TransactionScheduleFrequencyType = 3
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY       TransactionScheduleFrequencyType = 4
	TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS TransactionScheduleFrequencyType = 5
)

// TransactionTemplate represents transaction template stored in database
type TransactionTemplate struct {
	TemplateId                 int64                            `xorm:"PK"`
	Uid                        int64                            `xorm:"INDEX(IDX_transaction_template_uid_deleted_template_type_order) NOT NULL"`
	Deleted                    bool                             `xorm:"INDEX(IDX_transaction_template_uid_deleted_template_type_order) INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_time) NOT NULL"`
	TemplateType               TransactionTemplateType          `xorm:"INDEX(IDX_transaction_template_uid_deleted_template_type_order) INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_time) NOT NULL"`
	Name                       string                           `xorm:"VARCHAR(64) NOT NULL"`
	Type                       TransactionType                  `xorm:"NOT NULL"`
	CategoryId                 int64                            `xorm:"NOT NULL"`
	AccountId                  int64                            `xorm:"NOT NULL"`
	ScheduledFrequencyType     TransactionScheduleFrequencyType `xorm:"INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_time)"`
	ScheduledFrequency         string                           `xorm:"VARCHAR(100)"`
	ScheduledStartTime         *int64                           `xorm:"INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_time)"`
	ScheduledEndTime           *int64                           `xorm:"INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_time)"`
	ScheduledAt                int16                            `xorm:"INDEX(IDX_transaction_template_deleted_type_freqtype_scheduled_time)"`
	ScheduledTimezoneUtcOffset int16
	ScheduledTimezoneName      string `xorm:"VARCHAR(100)"`
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
	Name                       string                            `json:"name" binding:"required,notBlank,max=64"`
	Type                       TransactionType                   `json:"type" binding:"required"`
	CategoryId                 int64                             `json:"categoryId,string" binding:"required,min=1"`
	SourceAccountId            int64                             `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId       int64                             `json:"destinationAccountId,string" binding:"min=0"`
	SourceAmount               int64                             `json:"sourceAmount" binding:"validTransactionAmount"`
	DestinationAmount          int64                             `json:"destinationAmount" binding:"validTransactionAmount"`
	HideAmount                 bool                              `json:"hideAmount"`
	TagIds                     []string                          `json:"tagIds"`
	Comment                    string                            `json:"comment" binding:"max=255"`
	ScheduledFrequencyType     *TransactionScheduleFrequencyType `json:"scheduledFrequencyType" binding:"omitempty"`
	ScheduledFrequency         *string                           `json:"scheduledFrequency" binding:"omitempty"`
	ScheduledStartDate         *string                           `json:"scheduledStartDate" binding:"omitempty"`
	ScheduledEndDate           *string                           `json:"scheduledEndDate" binding:"omitempty"`
	ScheduledTimezoneUtcOffset *int16                            `json:"utcOffset" binding:"omitempty,min=-720,max=840"`
	ScheduledTimezoneName      *string                           `json:"timeZone" binding:"omitempty,max=100"`
	ClientSessionId            string                            `json:"clientSessionId"`
}

// TransactionTemplateModifyNameRequest represents all parameters of transaction template name modification request
type TransactionTemplateModifyNameRequest struct {
	Id   int64  `json:"id,string" binding:"required,min=1"`
	Name string `json:"name" binding:"required,notBlank,max=64"`
}

// TransactionTemplateModifyRequest represents all parameters of transaction template modification request
type TransactionTemplateModifyRequest struct {
	Id                         int64                             `json:"id,string" binding:"required,min=1"`
	Name                       string                            `json:"name" binding:"required,notBlank,max=64"`
	Type                       TransactionType                   `json:"type" binding:"required"`
	CategoryId                 int64                             `json:"categoryId,string" binding:"required,min=1"`
	SourceAccountId            int64                             `json:"sourceAccountId,string" binding:"required,min=1"`
	DestinationAccountId       int64                             `json:"destinationAccountId,string" binding:"min=0"`
	SourceAmount               int64                             `json:"sourceAmount" binding:"validTransactionAmount"`
	DestinationAmount          int64                             `json:"destinationAmount" binding:"validTransactionAmount"`
	HideAmount                 bool                              `json:"hideAmount"`
	TagIds                     []string                          `json:"tagIds"`
	Comment                    string                            `json:"comment" binding:"max=255"`
	ScheduledFrequencyType     *TransactionScheduleFrequencyType `json:"scheduledFrequencyType" binding:"omitempty"`
	ScheduledFrequency         *string                           `json:"scheduledFrequency" binding:"omitempty"`
	ScheduledStartDate         *string                           `json:"scheduledStartDate" binding:"omitempty"`
	ScheduledEndDate           *string                           `json:"scheduledEndDate" binding:"omitempty"`
	ScheduledTimezoneUtcOffset *int16                            `json:"utcOffset" binding:"omitempty,min=-720,max=840"`
	ScheduledTimezoneName      *string                           `json:"timeZone" binding:"omitempty,max=100"`
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
	ScheduledStartDate     *string                           `json:"scheduledStartDate" binding:"omitempty"`
	ScheduledEndDate       *string                           `json:"scheduledEndDate" binding:"omitempty"`
	ScheduledAt            *int16                            `json:"scheduledAt,omitempty"`
	ScheduledTimezoneName  *string                           `json:"timeZone,omitempty"`
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
		if t.ScheduledTimezoneName != "" {
			response.ScheduledTimezoneName = &t.ScheduledTimezoneName
		}

		templateTimeZone, err := t.GetScheduledTimezone()
		if err != nil {
			templateTimeZone = time.FixedZone("Template Timezone", int(t.ScheduledTimezoneUtcOffset)*60)
		}

		if t.ScheduledStartTime != nil {
			startDate := utils.FormatUnixTimeToLongDate(*t.ScheduledStartTime, templateTimeZone)
			response.ScheduledStartDate = &startDate
		}

		if t.ScheduledEndTime != nil {
			endDate := utils.FormatUnixTimeToLongDate(*t.ScheduledEndTime, templateTimeZone)
			response.ScheduledEndDate = &endDate
		}
	}

	return response
}

// GetScheduledTimezone returns the named zone or the saved fixed offset for older templates.
func (t *TransactionTemplate) GetScheduledTimezone() (*time.Location, error) {
	if t.ScheduledTimezoneName == "" {
		return time.FixedZone("Template Timezone", int(t.ScheduledTimezoneUtcOffset)*60), nil
	}
	if t.ScheduledTimezoneName == "Local" {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}
	return time.LoadLocation(t.ScheduledTimezoneName)
}

// ParseScheduledDate parses a schedule boundary in the template's time zone.
func (t *TransactionTemplate) ParseScheduledDate(date string, endOfDay bool) (time.Time, error) {
	if t.ScheduledTimezoneName == "" {
		if endOfDay {
			return utils.ParseFromLongDateLastTime(date, t.ScheduledTimezoneUtcOffset)
		}
		return utils.ParseFromLongDateFirstTime(date, t.ScheduledTimezoneUtcOffset)
	}
	location, err := t.GetScheduledTimezone()
	if err != nil {
		return time.Time{}, err
	}
	calendarDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	if endOfDay {
		calendarDate = calendarDate.AddDate(0, 0, 1)
	}
	boundary := time.Date(calendarDate.Year(), calendarDate.Month(), calendarDate.Day(), 0, 0, 0, 0, location)
	if boundary.Format("2006-01-02") != calendarDate.Format("2006-01-02") {
		// A clock change at midnight can move time.Date into the previous day.
		_, zoneEnd := boundary.ZoneBounds()
		if zoneEnd.IsZero() || zoneEnd.Format("2006-01-02") != calendarDate.Format("2006-01-02") {
			return time.Time{}, errs.ErrIncompleteOrIncorrectSubmission
		}
		boundary = zoneEnd
	}
	zoneStart, _ := boundary.ZoneBounds()
	if !zoneStart.IsZero() {
		// Use the first midnight when a backward clock change repeats it.
		_, previousOffset := zoneStart.Add(-time.Nanosecond).Zone()
		previousMidnight := time.Date(calendarDate.Year(), calendarDate.Month(), calendarDate.Day(), 0, 0, 0, 0, time.FixedZone("", previousOffset)).In(location)
		if previousMidnight.Before(boundary) && previousMidnight.Format("2006-01-02") == calendarDate.Format("2006-01-02") {
			boundary = previousMidnight
		}
	}
	if endOfDay {
		boundary = boundary.Add(-time.Nanosecond)
	}
	return boundary, nil
}

func (t *TransactionTemplate) toTransactionInfoResponse(utcOffset int16) *TransactionInfoResponse {
	tagIds := make([]string, 0, 0)

	if t.TagIds != "" {
		tagIds = strings.Split(t.TagIds, ",")
	}

	var destinationAmount *int64

	if t.Type == TRANSACTION_TYPE_TRANSFER {
		destinationAmount = &t.RelatedAccountAmount
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
		DestinationAmount:    destinationAmount,
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
