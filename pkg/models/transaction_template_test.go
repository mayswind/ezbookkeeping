package models

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionTemplateGetTagIds(t *testing.T) {
	template := &TransactionTemplate{
		TagIds: "1,2,3",
	}

	expectedValue := []int64{1, 2, 3}
	assert.EqualValues(t, expectedValue, template.GetTagIds())
}

func TestTransactionTemplateGetScheduledStatus(t *testing.T) {
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED,
	}

	actualValue := template.GetScheduledStatus(100)
	assert.Equal(t, SCHEDULED_TRANSACTION_STATUS_DISABLED, actualValue)

	endTime := int64(99)
	template.ScheduledEndTime = &endTime
	actualValue = template.GetScheduledStatus(100)
	assert.Equal(t, SCHEDULED_TRANSACTION_STATUS_DISABLED, actualValue)

	template.ScheduledFrequencyType = TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY
	actualValue = template.GetScheduledStatus(100)
	assert.Equal(t, SCHEDULED_TRANSACTION_STATUS_EXPIRED, actualValue)

	endTime = 100
	actualValue = template.GetScheduledStatus(100)
	assert.Equal(t, SCHEDULED_TRANSACTION_STATUS_ACTIVE, actualValue)

	endTime = 101
	actualValue = template.GetScheduledStatus(100)
	assert.Equal(t, SCHEDULED_TRANSACTION_STATUS_ACTIVE, actualValue)

	template.ScheduledEndTime = nil
	actualValue = template.GetScheduledStatus(100)
	assert.Equal(t, SCHEDULED_TRANSACTION_STATUS_ACTIVE, actualValue)

	startTime := int64(200)
	template.ScheduledStartTime = &startTime
	actualValue = template.GetScheduledStatus(100)
	assert.Equal(t, SCHEDULED_TRANSACTION_STATUS_ACTIVE, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_Disabled(t *testing.T) {
	transactionTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED,
	}

	actualValue := template.IsScheduledOnDate(transactionTime, []int64{0})
	assert.False(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_Daily(t *testing.T) {
	transactionTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
	}

	actualValue := template.IsScheduledOnDate(transactionTime, nil)
	assert.True(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_Weekly(t *testing.T) {
	transactionTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY,
	}

	actualValue := template.IsScheduledOnDate(transactionTime, []int64{0, 3})
	assert.True(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, []int64{1, 3})
	assert.False(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, nil)
	assert.False(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_Monthly(t *testing.T) {
	transactionTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
	}

	actualValue := template.IsScheduledOnDate(transactionTime, []int64{1, 15})
	assert.True(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, []int64{2, 15})
	assert.False(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_MonthEnd(t *testing.T) {
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
	}

	transactionTime := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)
	actualValue := template.IsScheduledOnDate(transactionTime, []int64{-1})
	assert.True(t, actualValue)

	transactionTime = time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC)
	actualValue = template.IsScheduledOnDate(transactionTime, []int64{-1})
	assert.True(t, actualValue)

	transactionTime = time.Date(2028, 2, 28, 0, 0, 0, 0, time.UTC)
	actualValue = template.IsScheduledOnDate(transactionTime, []int64{-1})
	assert.False(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, []int64{-2})
	assert.True(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_Yearly(t *testing.T) {
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY,
	}

	transactionTime := time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC)
	actualValue := template.IsScheduledOnDate(transactionTime, []int64{101, 229})
	assert.True(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, []int64{228, 301})
	assert.False(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_EveryNDaysMissingStartTime(t *testing.T) {
	transactionTime := time.Date(2028, 2, 29, 0, 0, 0, 0, time.UTC)
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
	}

	actualValue := template.IsScheduledOnDate(transactionTime, []int64{3})
	assert.False(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_EveryNDays(t *testing.T) {
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
	}

	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template.ScheduledStartTime = &startTime
	transactionTime := time.Unix(startTime, 0)
	actualValue := template.IsScheduledOnDate(transactionTime, []int64{3})
	assert.True(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime.AddDate(0, 0, 3), []int64{3})
	assert.True(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime.AddDate(0, 0, 2), []int64{3})
	assert.False(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime.AddDate(0, 0, -3), []int64{3})
	assert.False(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_EveryNDaysInvalidFrequency(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledStartTime:     &startTime,
	}

	transactionTime := time.Unix(startTime, 0)

	actualValue := template.IsScheduledOnDate(transactionTime, nil)
	assert.False(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, []int64{3, 6})
	assert.False(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, []int64{0})
	assert.False(t, actualValue)

	actualValue = template.IsScheduledOnDate(transactionTime, []int64{-3})
	assert.False(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_EveryNDaysFixedTimezone(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledStartTime:     &startTime,
	}

	location := time.FixedZone("UTC+8", 8*60*60)
	transactionTime := time.Date(2026, 1, 4, 0, 0, 0, 0, location)
	actualValue := template.IsScheduledOnDate(transactionTime, []int64{3})
	assert.True(t, actualValue)
}

func TestTransactionTemplateIsScheduledOnDate_InvalidFrequencyType(t *testing.T) {
	location := time.FixedZone("UTC+8", 8*60*60)
	transactionTime := time.Date(2026, 1, 4, 0, 0, 0, 0, location)
	template := &TransactionTemplate{
		ScheduledFrequencyType: TransactionScheduleFrequencyType(255),
	}

	actualValue := template.IsScheduledOnDate(transactionTime, []int64{0, 1})
	assert.False(t, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_InactiveTemplate(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_NORMAL,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledAt:            9 * 60,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, []int64{}, actualValue)

	template.TemplateType = TRANSACTION_TEMPLATE_TYPE_SCHEDULE
	template.ScheduledFrequencyType = TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED
	actualValue = template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, []int64{}, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_InvalidFrequency(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledAt:            9 * 60,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	template.ScheduledFrequency = ""
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, []int64{}, actualValue)

	template.ScheduledFrequency = "invalid"
	actualValue = template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, []int64{}, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_Daily(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledAt:            9 * 60,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	expectedValue := []int64{
		time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC).Unix(),
		time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC).Unix(),
		time.Date(2026, 2, 3, 9, 0, 0, 0, time.UTC).Unix(),
	}
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, expectedValue, actualValue)

	expectedValue = []int64{time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC).Unix()}
	actualValue = template.GetScheduledTimes(startTime, endTime, true)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_InvalidTimeRange(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledAt:            9 * 60,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	actualValue := template.GetScheduledTimes(startTime, startTime, false)
	assert.Equal(t, []int64{}, actualValue)

	actualValue = template.GetScheduledTimes(endTime, startTime, false)
	assert.Equal(t, []int64{}, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_TimeRangeBoundaries(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledAt:            9 * 60,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	firstTime := time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC).Unix()
	secondTime := time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC).Unix()
	expectedValue := []int64{firstTime}
	actualValue := template.GetScheduledTimes(firstTime, secondTime, false)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_StartAndEndBoundaries(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledAt:            9 * 60,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	secondTime := time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC).Unix()

	scheduledStartTime := secondTime
	scheduledEndTime := time.Date(2026, 2, 3, 9, 0, 0, 0, time.UTC).Unix()
	template.ScheduledStartTime = &scheduledStartTime
	template.ScheduledEndTime = &scheduledEndTime
	expectedValue := []int64{secondTime, scheduledEndTime}
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_NoOccurrenceInTimeRange(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 4, 0, 0, 0, 0, time.UTC).Unix()
	scheduledStartTime := time.Date(2026, 2, 2, 9, 0, 0, 0, time.UTC).Unix()
	scheduledEndTime := time.Date(2026, 2, 3, 9, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledStartTime:     &scheduledStartTime,
		ScheduledEndTime:       &scheduledEndTime,
		ScheduledAt:            9 * 60,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	scheduledStartTime = endTime
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, []int64{}, actualValue)

	template.ScheduledStartTime = nil
	scheduledEndTime = startTime - 1
	actualValue = template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, []int64{}, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_MultipleMonthDays(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "1,15,-1",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	endTime := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := []int64{
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix(),
		time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC).Unix(),
		time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC).Unix(),
	}
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledFrequency = "31"
	actualValue = template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, []int64{}, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_FixedTimezone(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:               TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType:     TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "0",
		ScheduledAt:                16 * 60,
		ScheduledTimezoneUtcOffset: 8 * 60,
		CreatedUnixTime:            startTime,
	}

	expectedValue := []int64{time.Date(2026, 2, 1, 16, 0, 0, 0, time.UTC).Unix()}
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneName = "Invalid/Timezone"
	actualValue = template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_TimezoneNameOverridesOffset(t *testing.T) {
	startTime := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:               TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType:     TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "0",
		ScheduledTimezoneName:      "UTC",
		ScheduledTimezoneUtcOffset: 8 * 60,
		CreatedUnixTime:            startTime,
	}

	expectedValue := []int64{startTime}
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionTemplateGetScheduledTimes_DSTTransition(t *testing.T) {
	template := &TransactionTemplate{
		TemplateType:               TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType:     TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "0",
		ScheduledTimezoneName:      "America/New_York",
		ScheduledTimezoneUtcOffset: 8 * 60,
	}

	template.ScheduledAt = 5 * 60
	template.CreatedUnixTime = time.Date(2026, 1, 1, 5, 0, 0, 0, time.UTC).Unix()
	startTime := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC).Unix()
	endTime := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := []int64{
		time.Date(2026, 3, 8, 5, 0, 0, 0, time.UTC).Unix(),
		time.Date(2026, 3, 9, 4, 0, 0, 0, time.UTC).Unix(),
	}
	actualValue := template.GetScheduledTimes(startTime, endTime, false)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_Daily(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        now,
	}

	expectedValue := now
	actualValue := template.GetNextScheduledTime(now)
	assert.Equal(t, &expectedValue, actualValue)

	expectedValue = time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC).Unix()
	actualValue = template.GetNextScheduledTime(now + 1)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_FutureStart(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        now,
	}

	startTime := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC).Unix()
	template.ScheduledStartTime = &startTime
	expectedValue := startTime
	actualValue := template.GetNextScheduledTime(now)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EndBoundary(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "0",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        now,
	}

	endTime := now - 1
	template.ScheduledEndTime = &endTime
	actualValue := template.GetNextScheduledTime(now)
	assert.Nil(t, actualValue)

	endTime = now
	expectedValue := now
	actualValue = template.GetNextScheduledTime(now)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_Disabled(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED,
		ScheduledFrequency:     "0",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        now,
	}

	actualValue := template.GetNextScheduledTime(now)
	assert.Nil(t, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_Weekly(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        now,
	}

	expectedValue := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC).Unix()
	actualValue := template.GetNextScheduledTime(now)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_MonthEnd(t *testing.T) {
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "-1",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}

	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC).Unix()
	actualValue := template.GetNextScheduledTime(now)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_NoOccurrenceBeforeEnd(t *testing.T) {
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "31",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}

	endTime := time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC).Unix()
	template.ScheduledEndTime = &endTime
	actualValue := template.GetNextScheduledTime(now)
	assert.Nil(t, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_LeapDayAcrossCentury(t *testing.T) {
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY,
		ScheduledFrequency:     "229",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}

	now := time.Date(2096, 3, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := time.Date(2104, 2, 29, 0, 0, 0, 0, time.UTC).Unix()
	actualValue := template.GetNextScheduledTime(now)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EveryNDays(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "10",
		ScheduledStartTime:     &startTime,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	expectedValue := startTime
	actualValue := template.GetNextScheduledTime(startTime)
	assert.Equal(t, &expectedValue, actualValue)

	expectedValue = time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC).Unix()
	actualValue = template.GetNextScheduledTime(startTime + 1)
	assert.Equal(t, &expectedValue, actualValue)

	template.ScheduledAt = 12 * 60
	expectedValue = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC).Unix()
	actualValue = template.GetNextScheduledTime(startTime + 60*60)
	assert.Equal(t, &expectedValue, actualValue)

	expectedValue = time.Date(2026, 1, 21, 12, 0, 0, 0, time.UTC).Unix()
	actualValue = template.GetNextScheduledTime(time.Date(2026, 1, 11, 13, 0, 0, 0, time.UTC).Unix())
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EveryNDaysLongInterval(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "10000",
		ScheduledStartTime:     &startTime,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	expectedValue := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 10000).Unix()
	actualValue := template.GetNextScheduledTime(startTime + 1)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EveryNDaysEndBeforeNextOccurrence(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "10",
		ScheduledStartTime:     &startTime,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	endTime := time.Date(2026, 1, 5, 23, 59, 59, 0, time.UTC).Unix()
	template.ScheduledEndTime = &endTime
	actualValue := template.GetNextScheduledTime(startTime + 1)
	assert.Nil(t, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EveryNDaysYearOutOfRange(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "4000000",
		ScheduledStartTime:     &startTime,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	actualValue := template.GetNextScheduledTime(startTime + 1)
	assert.Nil(t, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EveryNDaysInvalidFrequency(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "invalid",
		ScheduledStartTime:     &startTime,
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	actualValue := template.GetNextScheduledTime(startTime)
	assert.Nil(t, actualValue)

	template.ScheduledFrequency = "0"
	actualValue = template.GetNextScheduledTime(startTime)
	assert.Nil(t, actualValue)

	template.ScheduledFrequency = "-10"
	actualValue = template.GetNextScheduledTime(startTime)
	assert.Nil(t, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EveryNDaysMissingStartTime(t *testing.T) {
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "10",
		ScheduledTimezoneName:  "UTC",
		CreatedUnixTime:        startTime,
	}

	actualValue := template.GetNextScheduledTime(startTime)
	assert.Nil(t, actualValue)
}

func TestTransactionTemplateGetNextScheduledTime_EveryNDaysFixedTimezone(t *testing.T) {
	template := &TransactionTemplate{
		TemplateType:           TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledFrequencyType: TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "10",
		ScheduledTimezoneName:  "Invalid/Timezone",
	}

	location := time.FixedZone("UTC+8", 8*60*60)
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, location).Unix()
	template.ScheduledStartTime = &startTime
	template.ScheduledTimezoneUtcOffset = 8 * 60
	template.ScheduledAt = 16 * 60
	template.CreatedUnixTime = startTime
	expectedValue := time.Date(2026, 1, 11, 0, 0, 0, 0, location).Unix()
	actualValue := template.GetNextScheduledTime(startTime + 1)
	assert.Equal(t, &expectedValue, actualValue)
}

func TestTransactionTemplateInfoResponseSliceLess(t *testing.T) {
	var transactionTemplateRespSlice TransactionTemplateInfoResponseSlice
	transactionTemplateRespSlice = append(transactionTemplateRespSlice, &TransactionTemplateInfoResponse{
		TransactionInfoResponse: &TransactionInfoResponse{
			Id: 1,
		},
		DisplayOrder: 3,
	})
	transactionTemplateRespSlice = append(transactionTemplateRespSlice, &TransactionTemplateInfoResponse{
		TransactionInfoResponse: &TransactionInfoResponse{
			Id: 2,
		},
		DisplayOrder: 1,
	})
	transactionTemplateRespSlice = append(transactionTemplateRespSlice, &TransactionTemplateInfoResponse{
		TransactionInfoResponse: &TransactionInfoResponse{
			Id: 3,
		},
		DisplayOrder: 2,
	})

	sort.Sort(transactionTemplateRespSlice)

	assert.Equal(t, int64(2), transactionTemplateRespSlice[0].Id)
	assert.Equal(t, int64(3), transactionTemplateRespSlice[1].Id)
	assert.Equal(t, int64(1), transactionTemplateRespSlice[2].Id)
}

func TestScheduledTransactionOccurrenceResponseSliceLess(t *testing.T) {
	var occurrenceRespSlice ScheduledTransactionOccurrenceResponseSlice
	assert.Equal(t, 0, occurrenceRespSlice.Len())

	occurrenceRespSlice = append(occurrenceRespSlice, &ScheduledTransactionOccurrenceResponse{
		TemplateId:    3,
		ScheduledTime: 200,
	})
	occurrenceRespSlice = append(occurrenceRespSlice, &ScheduledTransactionOccurrenceResponse{
		TemplateId:    2,
		ScheduledTime: 100,
	})
	occurrenceRespSlice = append(occurrenceRespSlice, &ScheduledTransactionOccurrenceResponse{
		TemplateId:    1,
		ScheduledTime: 100,
	})

	assert.Equal(t, 3, occurrenceRespSlice.Len())
	assert.False(t, occurrenceRespSlice.Less(0, 0))
	assert.False(t, occurrenceRespSlice.Less(0, 1))
	assert.True(t, occurrenceRespSlice.Less(1, 0))
	assert.True(t, occurrenceRespSlice.Less(2, 1))

	sort.Sort(occurrenceRespSlice)

	assert.Equal(t, int64(1), occurrenceRespSlice[0].TemplateId)
	assert.Equal(t, int64(2), occurrenceRespSlice[1].TemplateId)
	assert.Equal(t, int64(3), occurrenceRespSlice[2].TemplateId)
}
