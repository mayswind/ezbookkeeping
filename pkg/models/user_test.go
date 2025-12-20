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
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(time.Now().Unix()), timezone))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsAll(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_ALL,
	}

	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(time.Now().Unix()), timezone))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsTodayOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_TODAY_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	todayFirstDatetime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	yesterdayLastDatetime := todayFirstDatetime.Add(-1 * time.Second)
	todayLastDatetime := yesterdayLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(todayFirstDatetime.Unix()), timezone))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(todayLastDatetime.Unix()), timezone))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(yesterdayLastDatetime.Unix()), timezone))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsLast24HourOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_LAST_24H_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	twentyfourHourBeforeDatetime := now.Add(-24 * time.Hour).Add(-1 * time.Second)

	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(twentyfourHourBeforeDatetime.Unix()), timezone))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(twentyfourHourBeforeDatetime.Add(1*time.Second).Unix()), timezone))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsThisWeekOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER,
		FirstDayOfWeek:       core.WEEKDAY_MONDAY,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	thisWeekFirstDatetime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if thisWeekFirstDatetime.Weekday() == time.Sunday {
		thisWeekFirstDatetime = thisWeekFirstDatetime.Add(-6 * 24 * time.Hour)
	} else if thisWeekFirstDatetime.Weekday() != time.Monday {
		thisWeekFirstDatetime = thisWeekFirstDatetime.Add(time.Duration(1-thisWeekFirstDatetime.Weekday()) * 24 * time.Hour)
	}

	lastWeekLastDatetime := thisWeekFirstDatetime.Add(-1 * time.Second)
	thisWeekLastDatetime := lastWeekLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisWeekFirstDatetime.Unix()), timezone))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisWeekLastDatetime.Unix()), timezone))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(lastWeekLastDatetime.Unix()), timezone))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsThisMonthOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	thisMonthFirstDatetime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	lastMonthLastDatetime := thisMonthFirstDatetime.Add(-1 * time.Second)
	thisMonthLastDatetime := lastMonthLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisMonthFirstDatetime.Unix()), timezone))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisMonthLastDatetime.Unix()), timezone))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(lastMonthLastDatetime.Unix()), timezone))
}

func TestUserCanEditTransactionByTransactionTime_ScopeIsThisYearOrLater(t *testing.T) {
	user := &User{
		TransactionEditScope: TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER,
	}

	now := time.Now()
	timezone := time.FixedZone("Timezone", int(utils.GetServerTimezoneOffsetMinutes())*60)
	thisYearFirstDatetime := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	lastYearLastDatetime := thisYearFirstDatetime.Add(-1 * time.Second)
	thisYearLastDatetime := lastYearLastDatetime.Add(24 * time.Hour)

	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisYearFirstDatetime.Unix()), timezone))
	assert.Equal(t, true, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(thisYearLastDatetime.Unix()), timezone))
	assert.Equal(t, false, user.CanEditTransactionByTransactionTime(utils.GetMinTransactionTimeFromUnixTime(lastYearLastDatetime.Unix()), timezone))
}
