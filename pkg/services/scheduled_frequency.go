package services

import (
	"errors"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// secondsPerDayInUtc is the length of a UTC day, which never has daylight saving transitions
const secondsPerDayInUtc int64 = 24 * 60 * 60

// errInvalidScheduledFrequencyType represents that the scheduled frequency type of a transaction
// template is not one of the supported ones, or that the template has no scheduled frequency at all
var errInvalidScheduledFrequencyType = errors.New("invalid scheduled transaction frequency")

// scheduledFrequencyMatchResult represents whether a transaction template fires at a given instant,
// and if it does not, why. The reasons exist so that callers can keep reporting the same detail they
// did when this logic was inlined in CreateScheduledTransactions.
type scheduledFrequencyMatchResult byte

const (
	scheduledFrequencyMatched scheduledFrequencyMatchResult = iota
	scheduledFrequencyUnmatchedWeekday
	scheduledFrequencyUnmatchedDayOfMonth
	scheduledFrequencyUnmatchedDayOfYear
	scheduledFrequencyInvalidEveryNDays
	scheduledFrequencyUnmatchedEveryNDays
	scheduledFrequencyBeforeStartTime
	scheduledFrequencyAfterEndTime
)

// GetScheduledOccurrences returns every instant at which the specified transaction template would
// create a transaction within the (from, to] range -- from is exclusive, to is inclusive. The
// returned instants are expressed in the timezone of the template.
//
// The occurrence of a given UTC day is that day's midnight in UTC plus ScheduledAt minutes, which is
// exactly how CreateScheduledTransactions builds the transaction time. Note that this instant is not
// the same calendar day as the UTC day it was derived from once it is read in the timezone of the
// template: with ScheduledAt near midnight and a negative UTC offset it lands on the previous day.
// That is why this function cannot take plain calendar dates.
//
// A template whose scheduled frequency is disabled yields no occurrences and no error. A template
// whose frequency is malformed yields an error.
func GetScheduledOccurrences(template *models.TransactionTemplate, from time.Time, to time.Time) ([]time.Time, error) {
	occurrences := make([]time.Time, 0)

	if template == nil {
		return nil, errInvalidScheduledFrequencyType
	}

	if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED {
		return occurrences, nil
	}

	frequencyValues, err := parseScheduledFrequencyValues(template)

	if err != nil {
		return nil, err
	}

	if !to.After(from) {
		return occurrences, nil
	}

	fromUnixTime := from.Unix()
	toUnixTime := to.Unix()
	templateTimeZone := getScheduledTemplateTimeZone(template)

	// A UTC day can only produce an occurrence inside [dayStart, dayStart+1day), so days before the
	// one containing "from" can never reach past it, and days after the one containing "to" always
	// overshoot it.
	firstDayUnixTime := getUtcDayStartUnixTime(fromUnixTime)
	lastDayUnixTime := getUtcDayStartUnixTime(toUnixTime)

	for dayUnixTime := firstDayUnixTime; dayUnixTime <= lastDayUnixTime; dayUnixTime += secondsPerDayInUtc {
		occurrenceUnixTime := dayUnixTime + int64(template.ScheduledAt)*60

		if occurrenceUnixTime <= fromUnixTime || occurrenceUnixTime > toUnixTime {
			continue
		}

		if matchScheduledFrequency(template, frequencyValues, occurrenceUnixTime) != scheduledFrequencyMatched {
			continue
		}

		occurrences = append(occurrences, time.Unix(occurrenceUnixTime, 0).In(templateTimeZone))
	}

	return occurrences, nil
}

// parseScheduledFrequencyValues validates the scheduled frequency type of the specified transaction
// template and parses its scheduled frequency into the list of values that the frequency type
// expects. The values are returned as stored, so the negative days of a monthly frequency are still
// negative and must be resolved against the month of each evaluated occurrence.
func parseScheduledFrequencyValues(template *models.TransactionTemplate) ([]int64, error) {
	if (template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY &&
		template.ScheduledFrequencyType != models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS) ||
		template.ScheduledFrequency == "" {
		return nil, errInvalidScheduledFrequencyType
	}

	return utils.StringArrayToInt64Array(strings.Split(template.ScheduledFrequency, ","))
}

// matchScheduledFrequency reports whether the specified transaction template would create a
// transaction at occurrenceUnixTime, which must be the instant the transaction would be created at
// (a UTC day start plus ScheduledAt minutes). frequencyValues must come from
// parseScheduledFrequencyValues.
//
// The daily frequency matches every day, so it has no check of its own.
func matchScheduledFrequency(template *models.TransactionTemplate, frequencyValues []int64, occurrenceUnixTime int64) scheduledFrequencyMatchResult {
	occurrenceTime := time.Unix(occurrenceUnixTime, 0).In(getScheduledTemplateTimeZone(template))

	if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY {
		if !containsInt64(frequencyValues, int64(occurrenceTime.Weekday())) {
			return scheduledFrequencyUnmatchedWeekday
		}
	} else if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY {
		if !matchesScheduledDayOfMonth(frequencyValues, occurrenceTime) {
			return scheduledFrequencyUnmatchedDayOfMonth
		}
	} else if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY {
		if !containsInt64(frequencyValues, int64(occurrenceTime.Month())*100+int64(occurrenceTime.Day())) {
			return scheduledFrequencyUnmatchedDayOfYear
		}
	} else if template.ScheduledFrequencyType == models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS {
		if template.ScheduledStartTime == nil || len(frequencyValues) != 1 || frequencyValues[0] <= 0 {
			return scheduledFrequencyInvalidEveryNDays
		}

		daysDiff := getScheduledEveryNDaysDiff(template, occurrenceUnixTime)

		if daysDiff < 0 || int64(daysDiff)%frequencyValues[0] != 0 {
			return scheduledFrequencyUnmatchedEveryNDays
		}
	}

	if template.ScheduledStartTime != nil && *template.ScheduledStartTime > occurrenceUnixTime {
		return scheduledFrequencyBeforeStartTime
	}

	if template.ScheduledEndTime != nil && *template.ScheduledEndTime < occurrenceUnixTime {
		return scheduledFrequencyAfterEndTime
	}

	return scheduledFrequencyMatched
}

// matchesScheduledDayOfMonth reports whether the day of month of the specified occurrence is one of
// the days of a monthly scheduled frequency. A negative value counts backwards from the end of the
// month, so -1 is the last day of the month.
//
// The length of the month is taken from the occurrence itself, read in the timezone of the template.
// Resolving it against any other month -- for instance the current month of the server clock, as
// this logic used to do while it was inlined in CreateScheduledTransactions -- returns the wrong day
// whenever the two months have different lengths.
//
// A positive day beyond the end of the month simply does not match: the occurrence is skipped for
// that month rather than moved to the last day of the month.
func matchesScheduledDayOfMonth(frequencyValues []int64, occurrenceTime time.Time) bool {
	dayOfMonth := int64(occurrenceTime.Day())
	maxDayOfMonth := int64(utils.GetMaxDayOfMonth(occurrenceTime.Year(), occurrenceTime.Month()))

	for i := 0; i < len(frequencyValues); i++ {
		frequencyValue := frequencyValues[i]

		if frequencyValue < 0 {
			frequencyValue = maxDayOfMonth + frequencyValue + 1
		}

		if frequencyValue == dayOfMonth {
			return true
		}
	}

	return false
}

// getScheduledEveryNDaysDiff returns how many whole days separate the specified occurrence from the
// start date of the template, both read in the timezone of the template. It is negative when the
// occurrence happens before the start date.
//
// The caller must have checked that ScheduledStartTime is set.
func getScheduledEveryNDaysDiff(template *models.TransactionTemplate, occurrenceUnixTime int64) int {
	templateTimeZone := getScheduledTemplateTimeZone(template)

	startTime := time.Unix(*template.ScheduledStartTime, 0).In(templateTimeZone)
	startDateOnly := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, templateTimeZone)

	occurrenceTime := time.Unix(occurrenceUnixTime, 0).In(templateTimeZone)
	occurrenceDateOnly := time.Date(occurrenceTime.Year(), occurrenceTime.Month(), occurrenceTime.Day(), 0, 0, 0, 0, templateTimeZone)

	return int(occurrenceDateOnly.Sub(startDateOnly).Hours() / 24)
}

// getScheduledTemplateTimeZone returns the timezone in which the scheduled frequency of the
// specified transaction template must be evaluated
func getScheduledTemplateTimeZone(template *models.TransactionTemplate) *time.Location {
	return time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
}

// getUtcDayStartUnixTime returns the unix time of the midnight in UTC of the day containing the
// specified unix time
func getUtcDayStartUnixTime(unixTime int64) int64 {
	dayTime := time.Unix(unixTime, 0).In(time.UTC)
	return time.Date(dayTime.Year(), dayTime.Month(), dayTime.Day(), 0, 0, 0, 0, time.UTC).Unix()
}

// containsInt64 reports whether the specified value is in the specified list. The scheduled
// frequency lists have at most a few dozen items, so a scan avoids building a set for every
// evaluated occurrence.
func containsInt64(values []int64, value int64) bool {
	for i := 0; i < len(values); i++ {
		if values[i] == value {
			return true
		}
	}

	return false
}
