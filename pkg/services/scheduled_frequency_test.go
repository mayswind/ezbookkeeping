package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

// This file is the phase 1 counterpart of transactions_scheduled_test.go: the same expectation table,
// now pointed at the real GetScheduledOccurrences engine instead of at a transcription of the logic
// that used to be inlined in CreateScheduledTransactions.
//
// See doc/financial-projections/Implementation-plan.md, phase 1.

// utcMidnight returns the midnight in UTC of the specified date
func utcMidnight(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// scheduledOccurrencesInUtcDays returns the occurrences of the specified template between the start
// of firstDay and the end of lastDay, both UTC days included, formatted in the timezone of the
// template so that a timezone shift is visible in the assertion.
func scheduledOccurrencesInUtcDays(t *testing.T, template *models.TransactionTemplate, firstDay time.Time, lastDay time.Time) []string {
	t.Helper()

	occurrences, err := GetScheduledOccurrences(template, firstDay.Add(-time.Second), lastDay.Add(24*time.Hour-time.Second))
	assert.Nil(t, err)

	return formatScheduledOccurrences(occurrences)
}

func formatScheduledOccurrences(occurrences []time.Time) []string {
	formatted := make([]string, 0, len(occurrences))

	for i := 0; i < len(occurrences); i++ {
		formatted = append(formatted, occurrences[i].Format("2006-01-02 15:04"))
	}

	return formatted
}

// --- WEEKLY ---------------------------------------------------------------------------------

func TestGetScheduledOccurrences_Weekly_SingleDay(t *testing.T) {
	// 2026-08-17 is a Monday
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "1")

	assert.Equal(t, []string{
		"2026-08-03 00:00",
		"2026-08-10 00:00",
		"2026-08-17 00:00",
		"2026-08-24 00:00",
		"2026-08-31 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 1), utcMidnight(2026, 8, 31)))
}

func TestGetScheduledOccurrences_Weekly_MultipleDays(t *testing.T) {
	// Mondays and Wednesdays
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "1,3")

	assert.Equal(t, []string{
		"2026-08-17 00:00",
		"2026-08-19 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 22)))
}

// TestGetScheduledOccurrences_Weekly_NegativeTemplateTimezoneShiftsTheEvaluatedDay is why the engine
// takes instants rather than calendar dates: 2026-08-17T01:00Z is still Sunday the 16th in UTC-3, so
// a template scheduled for Mondays does not fire and one scheduled for Sundays does.
func TestGetScheduledOccurrences_Weekly_NegativeTemplateTimezoneShiftsTheEvaluatedDay(t *testing.T) {
	mondayTemplate := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "1")
	mondayTemplate.ScheduledAt = 60
	mondayTemplate.ScheduledTimezoneUtcOffset = -180

	sundayTemplate := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "0")
	sundayTemplate.ScheduledAt = 60
	sundayTemplate.ScheduledTimezoneUtcOffset = -180

	assert.Equal(t, []string{}, scheduledOccurrencesInUtcDays(t, mondayTemplate, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 17)))
	assert.Equal(t, []string{"2026-08-16 22:00"}, scheduledOccurrencesInUtcDays(t, sundayTemplate, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 17)))
}

// TestGetScheduledOccurrences_Monthly_PositiveTemplateTimezoneShiftsTheEvaluatedDay is the mirror
// case: 2026-08-17T23:00Z is already the 18th in UTC+5:30.
func TestGetScheduledOccurrences_Monthly_PositiveTemplateTimezoneShiftsTheEvaluatedDay(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "18")
	template.ScheduledAt = 1380 // 23:00 UTC
	template.ScheduledTimezoneUtcOffset = 330

	assert.Equal(t, []string{"2026-08-18 04:30"}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 17)))
	assert.Equal(t, []string{}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 18), utcMidnight(2026, 8, 18)))
}

// --- MONTHLY --------------------------------------------------------------------------------

func TestGetScheduledOccurrences_Monthly_RegularDay(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "15")

	assert.Equal(t, []string{
		"2026-08-15 00:00",
		"2026-09-15 00:00",
		"2026-10-15 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 1), utcMidnight(2026, 10, 31)))
}

// TestGetScheduledOccurrences_Monthly_Day31IsSkippedInShortMonths preserves the behaviour of the
// cron: there is no clamping to the last day of the month, the occurrence is simply lost.
func TestGetScheduledOccurrences_Monthly_Day31IsSkippedInShortMonths(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "31")

	assert.Equal(t, []string{
		"2026-01-31 00:00",
		"2026-03-31 00:00",
		"2026-05-31 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 1, 1), utcMidnight(2026, 6, 30)))
}

// TestGetScheduledOccurrences_Monthly_NegativeDayIsResolvedPerMonth is the regression test for the
// bug pinned by TestLegacyScheduledFrequency_Monthly_NegativeDayResolvedAgainstServerLocalMonth.
//
// The inlined logic resolved the negative day against the month of the server-local clock, so a
// single call could only ever be right for one month. The engine resolves it against the month of
// each evaluated occurrence, in the timezone of the template, so a single range spanning months of
// 28, 30 and 31 days lands on the last day of every one of them.
func TestGetScheduledOccurrences_Monthly_NegativeDayIsResolvedPerMonth(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "-1")

	assert.Equal(t, []string{
		"2026-01-31 00:00",
		"2026-02-28 00:00",
		"2026-03-31 00:00",
		"2026-04-30 00:00",
		"2026-05-31 00:00",
		"2026-06-30 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 1, 1), utcMidnight(2026, 6, 30)))
}

func TestGetScheduledOccurrences_Monthly_NegativeDayInLeapFebruary(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "-1")

	assert.Equal(t, []string{
		"2024-02-29 00:00",
		"2024-03-31 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2024, 2, 1), utcMidnight(2024, 3, 31)))
}

func TestGetScheduledOccurrences_Monthly_NegativeDaySecondToLast(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "-2")

	assert.Equal(t, []string{
		"2026-02-27 00:00",
		"2026-03-30 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 2, 1), utcMidnight(2026, 3, 31)))
}

func TestGetScheduledOccurrences_Monthly_MixedPositiveAndNegativeDays(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "1,-1")

	assert.Equal(t, []string{
		"2026-02-01 00:00",
		"2026-02-28 00:00",
		"2026-03-01 00:00",
		"2026-03-31 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 2, 1), utcMidnight(2026, 3, 31)))
}

// --- DAILY ----------------------------------------------------------------------------------

func TestGetScheduledOccurrences_Daily(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")

	assert.Equal(t, []string{
		"2026-02-27 00:00",
		"2026-02-28 00:00",
		"2026-03-01 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 2, 27), utcMidnight(2026, 3, 1)))
}

func TestGetScheduledOccurrences_Daily_RespectsScheduledAt(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	template.ScheduledAt = 570 // 09:30 UTC

	assert.Equal(t, []string{
		"2026-08-17 09:30",
		"2026-08-18 09:30",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 18)))
}

// --- YEARLY ---------------------------------------------------------------------------------

func TestGetScheduledOccurrences_Yearly_MonthAndDayKey(t *testing.T) {
	// 0817 -> 17th of August
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, "817")

	assert.Equal(t, []string{
		"2026-08-17 00:00",
		"2027-08-17 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 1, 1), utcMidnight(2027, 12, 31)))
}

func TestGetScheduledOccurrences_Yearly_MultipleDates(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, "101,1225")

	assert.Equal(t, []string{
		"2026-01-01 00:00",
		"2026-12-25 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 1, 1), utcMidnight(2026, 12, 31)))
}

// TestGetScheduledOccurrences_Yearly_February29 preserves that a 29-February template fires only in
// leap years, with no fallback to the 28th.
func TestGetScheduledOccurrences_Yearly_February29(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, "229")

	assert.Equal(t, []string{"2024-02-29 00:00"}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2024, 1, 1), utcMidnight(2024, 12, 31)))
	assert.Equal(t, []string{}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 1, 1), utcMidnight(2026, 12, 31)))
}

// --- EVERY_N_DAYS ---------------------------------------------------------------------------

func TestGetScheduledOccurrences_EveryNDays_AlignedWithStartDate(t *testing.T) {
	startTime := utcUnixTime(2026, 8, 1, 0, 0)
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "7")
	template.ScheduledStartTime = &startTime

	assert.Equal(t, []string{
		"2026-08-01 00:00",
		"2026-08-08 00:00",
		"2026-08-15 00:00",
		"2026-08-22 00:00",
		"2026-08-29 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 1), utcMidnight(2026, 8, 31)))
}

func TestGetScheduledOccurrences_EveryNDays_CrossesMonthBoundary(t *testing.T) {
	startTime := utcUnixTime(2026, 8, 1, 0, 0)
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "10")
	template.ScheduledStartTime = &startTime

	assert.Equal(t, []string{
		"2026-08-01 00:00",
		"2026-08-11 00:00",
		"2026-08-21 00:00",
		"2026-08-31 00:00",
		"2026-09-10 00:00",
		"2026-09-20 00:00",
		"2026-09-30 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 1), utcMidnight(2026, 9, 30)))
}

func TestGetScheduledOccurrences_EveryNDays_BeforeStartDate(t *testing.T) {
	startTime := utcUnixTime(2026, 8, 10, 0, 0)
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "3")
	template.ScheduledStartTime = &startTime

	assert.Equal(t, []string{}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 1), utcMidnight(2026, 8, 9)))
}

// TestGetScheduledOccurrences_EveryNDays_InvalidConfiguration preserves that a misconfigured every
// N days template yields nothing rather than an error: the frequency itself parses, so the engine
// reports it the same way the cron did, by not firing.
func TestGetScheduledOccurrences_EveryNDays_InvalidConfiguration(t *testing.T) {
	startTime := utcUnixTime(2026, 8, 1, 0, 0)

	noStartDate := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "7")
	assert.Equal(t, []string{}, scheduledOccurrencesInUtcDays(t, noStartDate, utcMidnight(2026, 8, 1), utcMidnight(2026, 8, 31)))

	multipleValues := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "7,14")
	multipleValues.ScheduledStartTime = &startTime
	assert.Equal(t, []string{}, scheduledOccurrencesInUtcDays(t, multipleValues, utcMidnight(2026, 8, 1), utcMidnight(2026, 8, 31)))

	zeroInterval := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, "0")
	zeroInterval.ScheduledStartTime = &startTime
	assert.Equal(t, []string{}, scheduledOccurrencesInUtcDays(t, zeroInterval, utcMidnight(2026, 8, 1), utcMidnight(2026, 8, 31)))
}

// --- start / end time -----------------------------------------------------------------------

func TestGetScheduledOccurrences_StartTimeIsInclusive(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	startTime := utcUnixTime(2026, 8, 17, 0, 0)
	template.ScheduledStartTime = &startTime

	assert.Equal(t, []string{
		"2026-08-17 00:00",
		"2026-08-18 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 15), utcMidnight(2026, 8, 18)))
}

func TestGetScheduledOccurrences_EndTimeIsInclusive(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	endTime := utcUnixTime(2026, 8, 17, 0, 0)
	template.ScheduledEndTime = &endTime

	assert.Equal(t, []string{
		"2026-08-16 00:00",
		"2026-08-17 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 16), utcMidnight(2026, 8, 20)))
}

// TestGetScheduledOccurrences_EndTimeCutsTheRangeShort covers the case the projection depends on: a
// template that stops partway through the requested period.
func TestGetScheduledOccurrences_EndTimeCutsTheRangeShort(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "15")
	endTime := utcUnixTime(2026, 10, 1, 0, 0)
	template.ScheduledEndTime = &endTime

	assert.Equal(t, []string{
		"2026-08-15 00:00",
		"2026-09-15 00:00",
	}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 1), utcMidnight(2026, 12, 31)))
}

// TestGetScheduledOccurrences_StartTimeIsComparedAgainstTheScheduledInstant preserves that the start
// time is compared against the instant the transaction would be created at, not against the start of
// the day: a template starting at noon does not fire on its own start date if it is scheduled for
// the morning.
func TestGetScheduledOccurrences_StartTimeIsComparedAgainstTheScheduledInstant(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	template.ScheduledAt = 480 // 08:00 UTC
	startTime := utcUnixTime(2026, 8, 17, 12, 0)
	template.ScheduledStartTime = &startTime

	assert.Equal(t, []string{"2026-08-18 08:00"}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 18)))
}

// --- range boundaries -----------------------------------------------------------------------

// TestGetScheduledOccurrences_RangeStartIsExclusiveAndEndIsInclusive pins the half-open contract the
// projection relies on to split real and simulated amounts at "now" without double counting.
func TestGetScheduledOccurrences_RangeStartIsExclusiveAndEndIsInclusive(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")

	occurrences, err := GetScheduledOccurrences(template, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 19))

	assert.Nil(t, err)
	assert.Equal(t, []string{
		"2026-08-18 00:00",
		"2026-08-19 00:00",
	}, formatScheduledOccurrences(occurrences))
}

func TestGetScheduledOccurrences_EmptyRange(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")

	occurrences, err := GetScheduledOccurrences(template, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 17))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(occurrences))

	occurrences, err = GetScheduledOccurrences(template, utcMidnight(2026, 8, 19), utcMidnight(2026, 8, 17))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(occurrences))
}

// --- discarded templates --------------------------------------------------------------------

func TestGetScheduledOccurrences_DisabledFrequencyTypeYieldsNothing(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED, "1")

	occurrences, err := GetScheduledOccurrences(template, utcMidnight(2026, 8, 1), utcMidnight(2026, 8, 31))

	assert.Nil(t, err)
	assert.Equal(t, 0, len(occurrences))
}

func TestGetScheduledOccurrences_InvalidFrequencyReturnsError(t *testing.T) {
	from := utcMidnight(2026, 8, 1)
	to := utcMidnight(2026, 8, 31)

	emptyFrequency := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "")
	_, err := GetScheduledOccurrences(emptyFrequency, from, to)
	assert.NotNil(t, err)

	unparsableFrequency := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, "monday")
	_, err = GetScheduledOccurrences(unparsableFrequency, from, to)
	assert.NotNil(t, err)

	unknownFrequencyType := newScheduledTemplateForTest(models.TransactionScheduleFrequencyType(99), "1")
	_, err = GetScheduledOccurrences(unknownFrequencyType, from, to)
	assert.NotNil(t, err)

	_, err = GetScheduledOccurrences(nil, from, to)
	assert.NotNil(t, err)
}

// TestGetScheduledOccurrences_HiddenTemplateStillFires preserves that hiding a template does not
// stop it from creating transactions, so the projection must not filter on Hidden either.
func TestGetScheduledOccurrences_HiddenTemplateStillFires(t *testing.T) {
	template := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	template.Hidden = true

	assert.Equal(t, []string{"2026-08-17 00:00"}, scheduledOccurrencesInUtcDays(t, template, utcMidnight(2026, 8, 17), utcMidnight(2026, 8, 17)))
}

// --- equivalence with the pre-refactor cron -------------------------------------------------

// TestGetScheduledOccurrences_MatchesTheInlinedCronLogic is the differential test backing phase 1.4:
// for every day of a full year it compares the engine against
// legacyShouldCreateScheduledTransaction, the frozen transcription of the logic that used to live
// inside CreateScheduledTransactions.
//
// The templates are all on UTC and the legacy function is called with the server clock on the day
// being processed, so the month the legacy code resolves negative days against is the month of the
// occurrence. In that configuration the two must agree on every day, negative days included, which
// is what this test asserts.
//
// They diverge only when those two months differ -- a server that is not on UTC, around a month
// boundary. That is the bug fixed in phase 1: it is pinned by
// TestLegacyScheduledFrequency_Monthly_NegativeDayResolvedAgainstServerLocalMonth and its fix is
// covered by TestGetScheduledOccurrences_Monthly_NegativeDayIsResolvedPerMonth.
func TestGetScheduledOccurrences_MatchesTheInlinedCronLogic(t *testing.T) {
	everyNDaysStartTime := utcUnixTime(2026, 2, 3, 0, 0)
	boundedStartTime := utcUnixTime(2026, 4, 10, 0, 0)
	boundedEndTime := utcUnixTime(2026, 9, 20, 0, 0)

	templates := make([]*models.TransactionTemplate, 0)

	for _, frequency := range []string{"1", "0,6", "1,2,3,4,5"} {
		templates = append(templates, newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, frequency))
	}

	for _, frequency := range []string{"1", "15", "28", "29", "30", "31", "1,15,31", "-1", "-2", "-28", "1,-1"} {
		templates = append(templates, newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, frequency))
	}

	templates = append(templates, newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0"))

	for _, frequency := range []string{"101", "229", "817", "1231", "101,701"} {
		templates = append(templates, newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, frequency))
	}

	for _, frequency := range []string{"1", "2", "7", "13", "45"} {
		everyNDays := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, frequency)
		everyNDays.ScheduledStartTime = &everyNDaysStartTime
		templates = append(templates, everyNDays)
	}

	bounded := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, "0")
	bounded.ScheduledStartTime = &boundedStartTime
	bounded.ScheduledEndTime = &boundedEndTime
	templates = append(templates, bounded)

	scheduledAtNoon := newScheduledTemplateForTest(models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, "15")
	scheduledAtNoon.ScheduledAt = 720
	templates = append(templates, scheduledAtNoon)

	for _, template := range templates {
		for day := utcMidnight(2026, 1, 1); day.Year() < 2027; day = day.Add(24 * time.Hour) {
			dayStartUnixTime := day.Unix()

			expected := legacyShouldCreateScheduledTransaction(template, day, dayStartUnixTime)

			occurrences, err := GetScheduledOccurrences(template, day.Add(-time.Second), day.Add(24*time.Hour-time.Second))
			assert.Nil(t, err)

			assert.Equal(t, expected, len(occurrences) == 1,
				"frequency type %d, frequency %q, day %s", template.ScheduledFrequencyType, template.ScheduledFrequency, day.Format("2006-01-02"))
		}
	}
}
