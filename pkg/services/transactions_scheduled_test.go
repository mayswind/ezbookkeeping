package services

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// This file characterizes the CURRENT scheduled-frequency behaviour of
// TransactionService.CreateScheduledTransactions (transactions.go:820-902) before that logic is
// extracted into a reusable engine (see doc/financial-projections/Implementation-plan.md, phase 1).
//
// CreateScheduledTransactions cannot be called directly from a test: it queries the user data
// databases and needs a core.Context, and this package has no database test harness. So
// legacyShouldCreateScheduledTransaction below is a verbatim transcription of the decision logic of
// that function, with only the logging, the counters and the transaction creation removed. The
// value of these tests is the expectation table: phase 1 repoints the same cases at
// GetScheduledOccurrences, and any behavioural drift shows up there.
//
// Two quirks of the current behaviour are pinned deliberately and MUST be preserved by the
// extracted engine:
//
//  1. MONTHLY does not clamp to the end of the month. A template scheduled on day 31 simply does
//     not fire in months shorter than 31 days -- it is not moved to the 28th/30th.
//  2. Hidden templates still fire. The cron query filters on deleted/template_type/frequency, never
//     on the hidden column, so hiding a template does not stop it from creating transactions.
//
// A third behaviour is pinned as a KNOWN BUG rather than as intended semantics; see
// TestLegacyScheduledFrequency_Monthly_NegativeDayResolvedAgainstServerLocalMonth.
//
// Out of scope here: the SQL pre-filter on scheduled_at / scheduled_start_time / scheduled_end_time
// (transactions.go:785-798), which selects which templates reach this logic at all.

// legacyShouldCreateScheduledTransaction reproduces transactions.go:820-902.
//
// currentTime is the server-local wall clock (time.Unix(currentUnixTime, 0)) and
// todayFirstUnixTimeInUTC is midnight UTC of the day the cron is processing, exactly as
// CreateScheduledTransactions computes them.
func legacyShouldCreateScheduledTransaction(template *models.TransactionTemplate, currentTime time.Time, todayFirstUnixTimeInUTC int64) bool {
	if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED {
		return false
	}

	if (template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS) ||
		template.ScheduledFrequency == "" {
		return false
	}

	frequencyValues, err := utils.StringArrayToInt64Array(strings.Split(template.ScheduledFrequency, ","))

	if err != nil {
		return false
	}

	if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY {
		maxDayInMonth := utils.GetMaxDayOfMonth(currentTime.Year(), currentTime.Month())

		for i := 0; i < len(frequencyValues); i++ {
			if frequencyValues[i] < 0 {
				frequencyValues[i] = int64(maxDayInMonth) + frequencyValues[i] + 1
			}
		}
	}

	frequencyValueSet := utils.ToSet(frequencyValues)
	templateTimeZone := time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	transactionUnixTime := todayFirstUnixTimeInUTC + int64(template.ScheduledAt)*60
	transactionTime := time.Unix(transactionUnixTime, 0).In(templateTimeZone)

	if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY && !frequencyValueSet[int64(transactionTime.Weekday())] {
		return false
	} else if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY && !frequencyValueSet[int64(transactionTime.Day())] {
		return false
	} else if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY && !frequencyValueSet[int64(transactionTime.Month())*100+int64(transactionTime.Day())] {
		return false
	} else if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS {
		if template.ScheduledStartTime == nil || len(frequencyValues) != 1 || frequencyValues[0] <= 0 {
			return false
		}

		n := frequencyValues[0]
		startDate := time.Unix(*template.ScheduledStartTime, 0).In(templateTimeZone)
		startDateOnly := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, templateTimeZone)
		transactionDateOnly := time.Date(transactionTime.Year(), transactionTime.Month(), transactionTime.Day(), 0, 0, 0, 0, templateTimeZone)
		daysDiff := int(transactionDateOnly.Sub(startDateOnly).Hours() / 24)

		if daysDiff < 0 || int64(daysDiff)%n != 0 {
			return false
		}
	}

	if template.ScheduledStartTime != nil && *template.ScheduledStartTime > transactionUnixTime {
		return false
	}

	if template.ScheduledEndTime != nil && *template.ScheduledEndTime < transactionUnixTime {
		return false
	}

	return true
}

// utcDayStartUnixTime returns midnight UTC of the given date, the todayFirstUnixTimeInUTC the cron
// derives from the current instant.
func utcDayStartUnixTime(year int, month time.Month, day int) int64 {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix()
}

func utcUnixTime(year int, month time.Month, day int, hour int, minute int) int64 {
	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC).Unix()
}

func newScheduledTemplateForTest(frequencyType models.TransactionScheduleFrequencyType, frequency string) *models.TransactionTemplate {
	return &models.TransactionTemplate{
		TemplateId:             1001,
		Uid:                    1,
		TemplateType:           models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		CategoryId:             2001,
		AccountId:              3001,
		Amount:                 10000,
		ScheduledFrequencyType: frequencyType,
		ScheduledFrequency:     frequency,
	}
}

// --- WEEKLY ---------------------------------------------------------------------------------

func TestLegacyScheduledFrequency_Weekly_SingleDay(t *testing.T) {
	// 2026-08-17 is a Monday (weekday 1), 2026-08-19 a Wednesday (weekday 3)
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "1")
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 19)))
}

func TestLegacyScheduledFrequency_Weekly_MultipleDays(t *testing.T) {
	// Mondays and Wednesdays
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "1,3")
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 19)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 22)))
}

// TestLegacyScheduledFrequency_Weekly_TemplateTimezoneShiftsTheEvaluatedDay pins the reason the
// extracted engine cannot take a plain calendar date: the frequency is matched against the instant
// (midnight UTC + ScheduledAt minutes) converted to the template timezone, which can land on the
// previous or next day.
func TestLegacyScheduledFrequency_Weekly_TemplateTimezoneShiftsTheEvaluatedDay(t *testing.T) {
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	// 2026-08-17T01:00Z is still Sunday 2026-08-16T22:00 in UTC-3
	mondayTemplate := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "1")
	mondayTemplate.ScheduledAt = 60
	mondayTemplate.ScheduledTimezoneUtcOffset = -180

	sundayTemplate := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "0")
	sundayTemplate.ScheduledAt = 60
	sundayTemplate.ScheduledTimezoneUtcOffset = -180

	assert.False(t, legacyShouldCreateScheduledTransaction(mondayTemplate, currentTime, utcDayStartUnixTime(2026, 8, 17)))
	assert.True(t, legacyShouldCreateScheduledTransaction(sundayTemplate, currentTime, utcDayStartUnixTime(2026, 8, 17)))
}

// --- MONTHLY --------------------------------------------------------------------------------

func TestLegacyScheduledFrequency_Monthly_RegularDay(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "15")
	currentTime := time.Date(2026, 8, 15, 3, 0, 0, 0, time.UTC)

	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 15)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 16)))
}

// TestLegacyScheduledFrequency_Monthly_Day31IsSkippedInShortMonths pins quirk 1: there is no
// clamping to the last day of the month, the occurrence is simply lost.
func TestLegacyScheduledFrequency_Monthly_Day31IsSkippedInShortMonths(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "31")

	// 31-day month: fires
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 31, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 31)))

	// 30-day month: never fires, not even on the 30th
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 4, 30, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 4, 30)))

	// February: never fires, not even on the 28th
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 2, 28, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 2, 28)))
}

func TestLegacyScheduledFrequency_Monthly_NegativeDayIsLastDayOfMonth(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "-1")

	// 28-day February
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 2, 28, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 2, 28)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 2, 27, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 2, 27)))

	// 29-day February (leap year)
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2024, 2, 29, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2024, 2, 29)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2024, 2, 28, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2024, 2, 28)))

	// 30-day month
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 4, 30, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 4, 30)))

	// 31-day month
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 31, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 31)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 30, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 30)))
}

func TestLegacyScheduledFrequency_Monthly_NegativeDaySecondToLast(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "-2")

	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 2, 27, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 2, 27)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 2, 28, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 2, 28)))
}

// TestLegacyScheduledFrequency_Monthly_NegativeDayResolvedAgainstServerLocalMonth pins a KNOWN BUG,
// not intended behaviour.
//
// The negative day is resolved with GetMaxDayOfMonth(currentTime...), where currentTime is the
// SERVER-LOCAL wall clock (transactions.go:846), instead of the month of the occurrence being
// evaluated in the template timezone. The two disagree whenever the server-local date and the UTC
// date fall in different months -- which happens on every month boundary for any server not on UTC.
//
// The consequence for projections is much worse than for the cron: projecting N months ahead
// resolves every future "last day of the month" against today's month length.
//
// Phase 1 MUST fix this. When it does, the expectation below flips to true and this test moves to
// the fixed engine's suite as a regression test.
func TestLegacyScheduledFrequency_Monthly_NegativeDayResolvedAgainstServerLocalMonth(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "-1")

	// Server local clock is already in March (e.g. a UTC+5 server just past midnight) while the day
	// being processed is still 2026-02-28 in UTC. -1 resolves against March (31 days) instead of
	// February (28 days), so the last day of February does not fire.
	serverLocalTimeInMarch := time.Date(2026, 3, 1, 2, 0, 0, 0, time.UTC)
	assert.False(t, legacyShouldCreateScheduledTransaction(template, serverLocalTimeInMarch, utcDayStartUnixTime(2026, 2, 28)))

	// Same processed day, server local clock still in February: fires as expected.
	serverLocalTimeInFebruary := time.Date(2026, 2, 28, 21, 0, 0, 0, time.UTC)
	assert.True(t, legacyShouldCreateScheduledTransaction(template, serverLocalTimeInFebruary, utcDayStartUnixTime(2026, 2, 28)))
}

// --- DAILY ----------------------------------------------------------------------------------

func TestLegacyScheduledFrequency_Daily_AlwaysFires(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")

	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 17)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 18, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 18)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 2, 28, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 2, 28)))
}

// --- YEARLY ---------------------------------------------------------------------------------

func TestLegacyScheduledFrequency_Yearly_MonthAndDayKey(t *testing.T) {
	// 0817 -> 17th of August
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, "817")
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 18)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 7, 17)))
}

func TestLegacyScheduledFrequency_Yearly_MultipleDates(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, "101,1225")
	currentTime := time.Date(2026, 1, 1, 3, 0, 0, 0, time.UTC)

	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 1, 1)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 12, 25)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 12, 24)))
}

// TestLegacyScheduledFrequency_Yearly_February29 pins that a 29-Feb template fires only in leap
// years -- in common years the date does not exist and the occurrence is lost, with no fallback to
// the 28th.
func TestLegacyScheduledFrequency_Yearly_February29(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, "229")

	// Leap year: fires on 29 February
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2024, 2, 29, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2024, 2, 29)))

	// Common year: 28 February does not fire, and there is no 29th to process
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 2, 28, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 2, 28)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 3, 1, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 3, 1)))
}

// --- EVERY_N_DAYS ---------------------------------------------------------------------------

func TestLegacyScheduledFrequency_EveryNDays_AlignedWithStartDate(t *testing.T) {
	startTime := utcUnixTime(2026, 8, 1, 0, 0)
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "7")
	template.ScheduledStartTime = &startTime
	currentTime := time.Date(2026, 8, 15, 3, 0, 0, 0, time.UTC)

	// Day 0 of the cycle
	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 1)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 8)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 15)))

	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 10)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 14)))
}

func TestLegacyScheduledFrequency_EveryNDays_CrossesMonthBoundary(t *testing.T) {
	startTime := utcUnixTime(2026, 8, 1, 0, 0)
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "10")
	template.ScheduledStartTime = &startTime
	currentTime := time.Date(2026, 9, 10, 3, 0, 0, 0, time.UTC)

	// 2026-08-01 + 40 days
	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 9, 10)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 9, 11)))
}

func TestLegacyScheduledFrequency_EveryNDays_BeforeStartDate(t *testing.T) {
	startTime := utcUnixTime(2026, 8, 10, 0, 0)
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "3")
	template.ScheduledStartTime = &startTime
	currentTime := time.Date(2026, 8, 7, 3, 0, 0, 0, time.UTC)

	// Negative days diff: never fires before the start date
	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 7)))
}

func TestLegacyScheduledFrequency_EveryNDays_InvalidConfiguration(t *testing.T) {
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)
	startTime := utcUnixTime(2026, 8, 1, 0, 0)

	// No start date: the cycle has no origin
	noStartDate := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "7")
	assert.False(t, legacyShouldCreateScheduledTransaction(noStartDate, currentTime, utcDayStartUnixTime(2026, 8, 17)))

	// More than one frequency value
	multipleValues := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "7,14")
	multipleValues.ScheduledStartTime = &startTime
	assert.False(t, legacyShouldCreateScheduledTransaction(multipleValues, currentTime, utcDayStartUnixTime(2026, 8, 15)))

	// Non-positive interval
	zeroInterval := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "0")
	zeroInterval.ScheduledStartTime = &startTime
	assert.False(t, legacyShouldCreateScheduledTransaction(zeroInterval, currentTime, utcDayStartUnixTime(2026, 8, 1)))
}

// --- start / end time -----------------------------------------------------------------------

func TestLegacyScheduledFrequency_StartTimeIsInclusive(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	startTime := utcUnixTime(2026, 8, 17, 0, 0)
	template.ScheduledStartTime = &startTime

	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 16, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 16)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 17)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 18, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 18)))
}

func TestLegacyScheduledFrequency_EndTimeIsInclusive(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	endTime := utcUnixTime(2026, 8, 17, 0, 0)
	template.ScheduledEndTime = &endTime

	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 16, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 16)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 17)))
	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 18, 3, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 18)))
}

// TestLegacyScheduledFrequency_StartTimeIsComparedAgainstTheScheduledInstant pins that the start
// time is compared against midnight UTC + ScheduledAt, not against the start of the day: a template
// starting at noon does not fire on its own start date if it is scheduled for the morning.
func TestLegacyScheduledFrequency_StartTimeIsComparedAgainstTheScheduledInstant(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	template.ScheduledAt = 480 // 08:00 UTC
	startTime := utcUnixTime(2026, 8, 17, 12, 0)
	template.ScheduledStartTime = &startTime

	assert.False(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 17, 8, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 17)))
	assert.True(t, legacyShouldCreateScheduledTransaction(template, time.Date(2026, 8, 18, 8, 0, 0, 0, time.UTC), utcDayStartUnixTime(2026, 8, 18)))
}

// --- discarded templates --------------------------------------------------------------------

func TestLegacyScheduledFrequency_DisabledFrequencyType(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED, "1")
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
}

func TestLegacyScheduledFrequency_EmptyFrequency(t *testing.T) {
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	for _, frequencyType := range []models.TransactionScheduleFrequencyType{
		models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY,
		models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY,
		models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
	} {
		template := newScheduledTemplateForTest(frequencyType, "")
		assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
	}
}

func TestLegacyScheduledFrequency_UnparsableFrequency(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "monday")
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
}

func TestLegacyScheduledFrequency_UnknownFrequencyType(t *testing.T) {
	template := newScheduledTemplateForTest(models.TransactionScheduleFrequencyType(99), "1")
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	assert.False(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
}

// TestLegacyScheduledFrequency_HiddenTemplateStillFires pins quirk 2: hiding a template does not
// stop it from creating transactions, so the projection engine must not filter on Hidden either.
func TestLegacyScheduledFrequency_HiddenTemplateStillFires(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	template.Hidden = true
	currentTime := time.Date(2026, 8, 17, 3, 0, 0, 0, time.UTC)

	assert.True(t, legacyShouldCreateScheduledTransaction(template, currentTime, utcDayStartUnixTime(2026, 8, 17)))
}
