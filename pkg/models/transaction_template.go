package models

import (
	"strings"
	"time"

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

// ScheduledTransactionStatus represents the current status of a scheduled transaction
type ScheduledTransactionStatus byte

const (
	SCHEDULED_TRANSACTION_STATUS_ACTIVE   ScheduledTransactionStatus = 1
	SCHEDULED_TRANSACTION_STATUS_DISABLED ScheduledTransactionStatus = 2
	SCHEDULED_TRANSACTION_STATUS_EXPIRED  ScheduledTransactionStatus = 3
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
	ScheduledTimezoneName      string                           `xorm:"VARCHAR(100)"`
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
	TemplateType    TransactionTemplateType    `form:"template_type"`
	ScheduledStatus ScheduledTransactionStatus `form:"scheduled_status" binding:"omitempty,min=1,max=3"`
}

// ScheduledTransactionMonthRequest represents all parameters of scheduled transaction month request
type ScheduledTransactionMonthRequest struct {
	Year  int `form:"year" binding:"required,min=1"`
	Month int `form:"month" binding:"required,min=1"`
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
	ScheduledTimezoneName      *string                           `json:"timeZone" binding:"omitempty,min=3,max=100"`
	ScheduledTimezoneUtcOffset *int16                            `json:"utcOffset" binding:"omitempty,min=-720,max=840"`
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
	ScheduledTimezoneName      *string                           `json:"timeZone" binding:"omitempty,min=3,max=100"`
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
	ScheduledStartDate     *string                           `json:"scheduledStartDate" binding:"omitempty"`
	ScheduledEndDate       *string                           `json:"scheduledEndDate" binding:"omitempty"`
	ScheduledAt            *int16                            `json:"scheduledAt,omitempty"`
	ScheduledStatus        ScheduledTransactionStatus        `json:"scheduledStatus,omitempty"`
	NextScheduledTime      *int64                            `json:"nextScheduledTime,omitempty"`
	ScheduledTimezoneName  string                            `json:"timeZone,omitempty"`
	DisplayOrder           int32                             `json:"displayOrder"`
	Hidden                 bool                              `json:"hidden"`
}

// ScheduledTransactionOccurrenceResponse represents a scheduled transaction occurrence
type ScheduledTransactionOccurrenceResponse struct {
	TemplateId    int64 `json:"templateId,string"`
	ScheduledTime int64 `json:"scheduledTime"`
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

// GetScheduledStatus returns the status of the scheduled transaction template at the given time
func (t *TransactionTemplate) GetScheduledStatus(now int64) ScheduledTransactionStatus {
	if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED {
		return SCHEDULED_TRANSACTION_STATUS_DISABLED
	} else if t.ScheduledEndTime != nil && *t.ScheduledEndTime < now {
		return SCHEDULED_TRANSACTION_STATUS_EXPIRED
	}

	return SCHEDULED_TRANSACTION_STATUS_ACTIVE
}

// IsScheduledOnDate returns whether the transaction template is scheduled on the given date
func (t *TransactionTemplate) IsScheduledOnDate(transactionTime time.Time, frequencyValues []int64) bool {
	if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED {
		return false
	} else if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY {
		return true
	}

	if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS {
		if t.ScheduledStartTime == nil || len(frequencyValues) != 1 || frequencyValues[0] <= 0 {
			return false
		}

		n := frequencyValues[0]
		startDate := time.Unix(*t.ScheduledStartTime, 0).In(transactionTime.Location())
		daysDiff := utils.GetDaysBetweenDates(startDate, transactionTime)
		return daysDiff >= 0 && int64(daysDiff)%n == 0
	}

	for _, value := range frequencyValues {
		if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY && value == int64(transactionTime.Weekday()) {
			return true
		} else if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY {
			if value < 0 {
				maxDayInMonth := utils.GetMaxDayOfMonth(transactionTime.Year(), transactionTime.Month())
				value = int64(maxDayInMonth) + value + 1
			}

			if value == int64(transactionTime.Day()) {
				return true
			}
		} else if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY && value == int64(transactionTime.Month())*100+int64(transactionTime.Day()) {
			return true
		}
	}

	return false
}

// GetScheduledTimes returns all planned times of the transaction template between the given time range
func (t *TransactionTemplate) GetScheduledTimes(startTime, endTime int64, firstOnly bool) []int64 {
	times := make([]int64, 0)

	if t.TemplateType != TRANSACTION_TEMPLATE_TYPE_SCHEDULE || t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED || t.ScheduledFrequency == "" {
		return times
	}

	frequencyValues, err := utils.StringArrayToInt64Array(strings.Split(t.ScheduledFrequency, ","))

	if err != nil {
		return times
	}

	if t.ScheduledStartTime != nil && *t.ScheduledStartTime > startTime {
		startTime = *t.ScheduledStartTime
	}

	if t.ScheduledEndTime != nil && *t.ScheduledEndTime < endTime-1 {
		endTime = *t.ScheduledEndTime + 1
	}

	if startTime >= endTime {
		return times
	}

	templateTimezone := time.FixedZone("Template Timezone", int(t.ScheduledTimezoneUtcOffset)*60)

	if len(t.ScheduledTimezoneName) > 0 {
		location, err := time.LoadLocation(t.ScheduledTimezoneName)

		if err == nil {
			templateTimezone = location
		}
	}

	startDate := time.Unix(startTime, 0).UTC()
	day := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -2)

	for ; day.Unix() < endTime; day = day.AddDate(0, 0, 1) {
		transactionTime := utils.GetScheduledTransactionFinalTime(day.Unix(), t.ScheduledAt, t.CreatedUnixTime, templateTimezone)

		if transactionTime == nil || transactionTime.Unix() < startTime || transactionTime.Unix() >= endTime || !t.IsScheduledOnDate(*transactionTime, frequencyValues) {
			continue
		}

		times = append(times, transactionTime.Unix())

		if firstOnly {
			break
		}
	}

	return times
}

// GetNextScheduledTime returns the next planned time of the transaction template after the given time
func (t *TransactionTemplate) GetNextScheduledTime(now int64) *int64 {
	if t.GetScheduledStatus(now) != SCHEDULED_TRANSACTION_STATUS_ACTIVE {
		return nil
	}

	startTime := now

	if t.ScheduledStartTime != nil && *t.ScheduledStartTime > startTime {
		startTime = *t.ScheduledStartTime
	}

	if t.ScheduledFrequencyType == TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS {
		frequency, err := utils.StringToInt(t.ScheduledFrequency)

		if err != nil || frequency <= 0 || t.ScheduledStartTime == nil {
			return nil
		}

		templateTimezone := time.FixedZone("Template Timezone", int(t.ScheduledTimezoneUtcOffset)*60)

		if len(t.ScheduledTimezoneName) > 0 {
			location, err := time.LoadLocation(t.ScheduledTimezoneName)

			if err == nil {
				templateTimezone = location
			}
		}

		anchor := time.Unix(*t.ScheduledStartTime, 0).In(templateTimezone)
		date := time.Unix(startTime, 0).In(templateTimezone)
		daysDiff := utils.GetDaysBetweenDates(anchor, date)

		days := daysDiff / frequency * frequency
		candidate := anchor.AddDate(0, 0, days)

		for attempt := 0; attempt < 2; attempt++ {
			times := t.GetScheduledTimes(candidate.AddDate(0, 0, -2).Unix(), candidate.AddDate(0, 0, 3).Unix(), false)

			for _, scheduledTime := range times {
				if scheduledTime >= startTime {
					return &scheduledTime
				}
			}

			if frequency > utils.GetDaysBetweenDates(candidate, time.Date(9999, 12, 31, 0, 0, 0, 0, templateTimezone)) {
				return nil
			}

			candidate = candidate.AddDate(0, 0, frequency)
		}

		return nil
	}

	// Eight years include the longest gap between Gregorian leap days
	endTime := time.Unix(startTime, 0).AddDate(8, 0, 1).Unix()
	times := t.GetScheduledTimes(startTime, endTime, true)

	if len(times) == 0 {
		return nil
	}

	return &times[0]
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
		now := time.Now().Unix()
		response.ScheduledStatus = t.GetScheduledStatus(now)
		response.NextScheduledTime = t.GetNextScheduledTime(now)
		response.ScheduledFrequencyType = &t.ScheduledFrequencyType
		response.ScheduledFrequency = &t.ScheduledFrequency
		response.ScheduledAt = &t.ScheduledAt
		response.ScheduledTimezoneName = t.ScheduledTimezoneName

		templateTimezone := time.FixedZone("Template Timezone", int(t.ScheduledTimezoneUtcOffset)*60)

		if len(t.ScheduledTimezoneName) > 0 {
			if location, err := time.LoadLocation(t.ScheduledTimezoneName); err == nil {
				templateTimezone = location
			}
		}

		if t.ScheduledStartTime != nil {
			startDate := utils.FormatUnixTimeToLongDate(*t.ScheduledStartTime, templateTimezone)
			response.ScheduledStartDate = &startDate
		}

		if t.ScheduledEndTime != nil {
			endDate := utils.FormatUnixTimeToLongDate(*t.ScheduledEndTime, templateTimezone)
			response.ScheduledEndDate = &endDate
		}
	}

	return response
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

// ScheduledTransactionOccurrenceResponseSlice represents the slice data structure of ScheduledTransactionOccurrenceResponse
type ScheduledTransactionOccurrenceResponseSlice []*ScheduledTransactionOccurrenceResponse

// Len returns the count of items
func (s ScheduledTransactionOccurrenceResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s ScheduledTransactionOccurrenceResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s ScheduledTransactionOccurrenceResponseSlice) Less(i, j int) bool {
	if s[i].ScheduledTime == s[j].ScheduledTime {
		return s[i].TemplateId < s[j].TemplateId
	}

	return s[i].ScheduledTime < s[j].ScheduledTime
}
