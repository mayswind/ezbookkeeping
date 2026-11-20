package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestUserCanEditTransactionByTransactionTime_ScopeIsNone(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_NONE,
	}

	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(time.Now().Unix()), timezone, nil, nil))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsAll(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_ALL,
	}

	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(time.Now().Unix()), timezone, nil, nil))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsTodayOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_TODAY_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	todayFirstDatetime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timezone)
	yesterdayLastDatetime := todayFirstDatetime.Add(-1 * time.Second)
	todayLastDatetime := yesterdayLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(todayFirstDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(todayFirstDatetime.Add(1*time.Second).Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(todayLastDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(yesterdayLastDatetime.Unix()), timezone, nil, nil))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsLast24HourOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_LAST_24H_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	twentyfourHourBeforeDatetime := now.Add(-24 * time.Hour).Add(-1 * time.Second)

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(twentyfourHourBeforeDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(twentyfourHourBeforeDatetime.Add(1*time.Second).Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(twentyfourHourBeforeDatetime.Add(2*time.Second).Unix()), timezone, nil, nil))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsThisWeekOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER,
		FirstDayOfWeek:       core.WEEKDAY_MONDAY,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	thisWeekFirstDatetime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timezone)

	if thisWeekFirstDatetime.Weekday() == time.Sunday {
		thisWeekFirstDatetime = thisWeekFirstDatetime.Add(-6 * 24 * time.Hour)
	} else if thisWeekFirstDatetime.Weekday() != time.Monday {
		thisWeekFirstDatetime = thisWeekFirstDatetime.Add(time.Duration(1-thisWeekFirstDatetime.Weekday()) * 24 * time.Hour)
	}

	lastWeekLastDatetime := thisWeekFirstDatetime.Add(-1 * time.Second)
	thisWeekLastDatetime := lastWeekLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisWeekFirstDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisWeekFirstDatetime.Add(1*time.Second).Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisWeekLastDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(lastWeekLastDatetime.Unix()), timezone, nil, nil))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsThisMonthOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	thisMonthFirstDatetime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, timezone)
	lastMonthLastDatetime := thisMonthFirstDatetime.Add(-1 * time.Second)
	thisMonthLastDatetime := lastMonthLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisMonthFirstDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisMonthFirstDatetime.Add(1*time.Second).Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisMonthLastDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(lastMonthLastDatetime.Unix()), timezone, nil, nil))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsThisYearOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	thisYearFirstDatetime := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, timezone)
	lastYearLastDatetime := thisYearFirstDatetime.Add(-1 * time.Second)
	thisYearLastDatetime := lastYearLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisYearFirstDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisYearFirstDatetime.Add(1*time.Second).Unix()), timezone, nil, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisYearLastDatetime.Unix()), timezone, nil, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(lastYearLastDatetime.Unix()), timezone, nil, nil))
}

func TestUserCanEditTransactionByTransactionTime_ScopeStartsOnClientDate(t *testing.T) {
	// At any time of day, at least one of these two zones is on a different date than the server
	timezones := []*time.Location{
		time.FixedZone("UTC+14", 14*60*60),
		time.FixedZone("UTC-12", -12*60*60),
	}

	for _, timezone := range timezones {
		clientNow := time.Now().In(timezone)
		clientToday := time.Date(clientNow.Year(), clientNow.Month(), clientNow.Day(), 0, 0, 0, 0, timezone)
		daysSinceMonday := (int(clientNow.Weekday()) + 6) % 7

		scopeStartTimes := map[TransactionEditScope]time.Time{
			TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER:  clientToday.AddDate(0, 0, -daysSinceMonday),
			TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER: time.Date(clientNow.Year(), clientNow.Month(), 1, 0, 0, 0, 0, timezone),
			TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER:  time.Date(clientNow.Year(), 1, 1, 0, 0, 0, 0, timezone),
		}

		for scope, scopeStartTime := range scopeStartTimes {
			user := &User{
				TransactionEditScope: scope,
				FirstDayOfWeek:       core.WEEKDAY_MONDAY,
			}

			assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(scopeStartTime.Add(1*time.Second).Unix()), timezone, nil, nil), "scope %d, %s", scope, scopeStartTime.Add(1*time.Second))
			assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(scopeStartTime.Add(-1*time.Second).Unix()), timezone, nil, nil), "scope %d, %s", scope, scopeStartTime.Add(-1*time.Second))
		}
	}
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsLastReconciledTimeOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope:  TRANSACTION_EDIT_SCOPE_LAST_RECONCILED_TIME_OR_LATER,
		UseLastReconciledTime: true,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	sourceAccountLastReconciledTime := now.Add(-24 * time.Hour)
	sourceAccountLastRecondiledUnixTime := sourceAccountLastReconciledTime.Unix()
	sourceAccount := &Account{
		Extend: &AccountExtend{
			LastReconciledTime: &sourceAccountLastRecondiledUnixTime,
		},
	}
	destinationAccountLastReconciledTime := now.Add(-20 * time.Hour)
	destinationAccountLastReconciledUnixTime := destinationAccountLastReconciledTime.Unix()
	destinationAccount := &Account{
		Extend: &AccountExtend{
			LastReconciledTime: &destinationAccountLastReconciledUnixTime,
		},
	}

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(sourceAccountLastReconciledTime.Add(-1*time.Second).Unix()), timezone, sourceAccount, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(sourceAccountLastReconciledTime.Unix()), timezone, sourceAccount, nil))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(sourceAccountLastReconciledTime.Add(1*time.Second).Unix()), timezone, sourceAccount, nil))

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(destinationAccountLastReconciledTime.Add(-1*time.Second).Unix()), timezone, sourceAccount, destinationAccount))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(destinationAccountLastReconciledTime.Unix()), timezone, sourceAccount, destinationAccount))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(destinationAccountLastReconciledTime.Add(1*time.Second).Unix()), timezone, sourceAccount, destinationAccount))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsLastReconciledTimeOrLaterButUserDoesNotUseLastReconciledTime(t *testing.T) {
	user := &User{
		TransactionEditScope:  TRANSACTION_EDIT_SCOPE_LAST_RECONCILED_TIME_OR_LATER,
		UseLastReconciledTime: false,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	sourceAccountLastReconciledTime := now.Add(-24 * time.Hour)
	sourceAccountLastRecondiledUnixTime := sourceAccountLastReconciledTime.Unix()
	sourceAccount := &Account{
		Extend: &AccountExtend{
			LastReconciledTime: &sourceAccountLastRecondiledUnixTime,
		},
	}
	destinationAccountLastReconciledTime := now.Add(-20 * time.Hour)
	destinationAccountLastReconciledUnixTime := destinationAccountLastReconciledTime.Unix()
	destinationAccount := &Account{
		Extend: &AccountExtend{
			LastReconciledTime: &destinationAccountLastReconciledUnixTime,
		},
	}

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(sourceAccountLastReconciledTime.Add(-1*time.Second).Unix()), timezone, sourceAccount, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(sourceAccountLastReconciledTime.Unix()), timezone, sourceAccount, nil))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(sourceAccountLastReconciledTime.Add(1*time.Second).Unix()), timezone, sourceAccount, nil))

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(destinationAccountLastReconciledTime.Add(-1*time.Second).Unix()), timezone, sourceAccount, destinationAccount))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(destinationAccountLastReconciledTime.Unix()), timezone, sourceAccount, destinationAccount))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(destinationAccountLastReconciledTime.Add(1*time.Second).Unix()), timezone, sourceAccount, destinationAccount))
}
