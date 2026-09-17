package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func TestGetScheduledTransactionTimeAndFrequencyValues_LastDayOfMonthInTemplateTimezone(t *testing.T) {
	// The server runs east of the template time zone, far enough that it has already
	// entered the next month when the template is due.
	original := time.Local
	time.Local = time.FixedZone("Server Timezone", 14*60*60)
	defer func() { time.Local = original }()

	// A template in UTC-12 scheduled for the last day of each month. It is due at
	// 2027-02-28 12:00 UTC, which is 2027-02-28 00:00 in UTC-12 and 2027-03-01 02:00
	// for the server.
	template := &models.TransactionTemplate{
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledTimezoneUtcOffset: -720,
		ScheduledAt:                720,
	}
	todayFirstUnixTimeInUTC := time.Date(2027, 2, 28, 0, 0, 0, 0, time.UTC).Unix()

	transactionTime, frequencyValues := getScheduledTransactionTimeAndFrequencyValues(template, []int64{-1}, todayFirstUnixTimeInUTC)

	assert.Equal(t, 2, int(transactionTime.Month()))
	assert.Equal(t, 28, transactionTime.Day())
	// February has 28 days in 2027, so -1 must resolve to 28, not to the 31 days of
	// the month the server happens to be in.
	assert.Equal(t, []int64{28}, frequencyValues)
}

func TestGetScheduledTransactionTimeAndFrequencyValues_PositiveValuesUnchanged(t *testing.T) {
	template := &models.TransactionTemplate{
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledTimezoneUtcOffset: 180,
		ScheduledAt:                1260,
	}
	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	transactionTime, frequencyValues := getScheduledTransactionTimeAndFrequencyValues(template, []int64{1}, todayFirstUnixTimeInUTC)

	assert.Equal(t, 1, transactionTime.Day())
	assert.Equal(t, 12, int(transactionTime.Month()))
	assert.Equal(t, []int64{1}, frequencyValues)
}

func TestGetScheduledTransactionTimeAndFrequencyValues_TemplateMonthAheadOfUtcMonth(t *testing.T) {
	// A template in UTC+3 due on 2027-05-01 00:00 local time runs at 2027-04-30 21:00
	// UTC, so the template month and the UTC month differ. -31 has to be resolved
	// against May, which has 31 days, and therefore yields the first of the month.
	original := time.Local
	time.Local = time.UTC
	defer func() { time.Local = original }()

	template := &models.TransactionTemplate{
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledTimezoneUtcOffset: 180,
		ScheduledAt:                1260,
	}
	todayFirstUnixTimeInUTC := time.Date(2027, 4, 30, 0, 0, 0, 0, time.UTC).Unix()

	transactionTime, frequencyValues := getScheduledTransactionTimeAndFrequencyValues(template, []int64{-31}, todayFirstUnixTimeInUTC)

	assert.Equal(t, 5, int(transactionTime.Month()))
	assert.Equal(t, 1, transactionTime.Day())
	assert.Equal(t, []int64{1}, frequencyValues)
}
