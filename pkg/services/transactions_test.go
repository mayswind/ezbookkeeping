package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

func TestCreateTransactionByScheduledTransactionTemplate_DisabledFrequency(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DISABLED,
		ScheduledFrequency:     "1",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_InvalidFrequencyType(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TransactionScheduleFrequencyType(255),
		ScheduledFrequency:     "1",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_EmptyFrequency(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequencyType = models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequencyType = models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequencyType = models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequencyType = models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_InvalidFrequency(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "invalid",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequency = "1,invalid"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequency = "9223372036854775808"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequency = "1,"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_Daily(t *testing.T) {
	template := &models.TransactionTemplate{
		TemplateId:             1001,
		Uid:                    123,
		TemplateType:           models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		Name:                   "Scheduled Expense",
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		CategoryId:             234,
		AccountId:              345,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		TagIds:                 "1001,1002,1003",
		Amount:                 456,
		RelatedAccountId:       567,
		RelatedAccountAmount:   678,
		HideAmount:             true,
		Comment:                "Scheduled transaction comment",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	expectedValue := &models.Transaction{
		Uid:              123,
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		CategoryId:       234,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		AccountId:        345,
		Amount:           456,
		HideAmount:       true,
		Comment:          "Scheduled transaction comment",
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	expectedTagIds := []int64{1001, 1002, 1003}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, expectedTagIds, actualTagIds)

	// The following day creates another transaction with the same template data.
	startTimeInUTC = time.Date(2026, 12, 2, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, expectedTagIds, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_Income(t *testing.T) {
	template := &models.TransactionTemplate{
		Uid:                    123,
		Type:                   models.TRANSACTION_TYPE_INCOME,
		CategoryId:             234,
		AccountId:              345,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		Amount:                 456,
		RelatedAccountId:       567,
		RelatedAccountAmount:   678,
		Comment:                "Scheduled income",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	expectedValue := &models.Transaction{
		Uid:              123,
		Type:             models.TRANSACTION_DB_TYPE_INCOME,
		CategoryId:       234,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		AccountId:        345,
		Amount:           456,
		Comment:          "Scheduled income",
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_Transfer(t *testing.T) {
	template := &models.TransactionTemplate{
		Uid:                    123,
		Type:                   models.TRANSACTION_TYPE_TRANSFER,
		CategoryId:             234,
		AccountId:              345,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		TagIds:                 "1001,1002",
		Amount:                 456,
		RelatedAccountId:       567,
		RelatedAccountAmount:   678,
		HideAmount:             true,
		Comment:                "Scheduled transfer",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	expectedValue := &models.Transaction{
		Uid:                  123,
		Type:                 models.TRANSACTION_DB_TYPE_TRANSFER_OUT,
		CategoryId:           234,
		TransactionTime:      startTimeInUTC.Unix() * 1000,
		AccountId:            345,
		Amount:               456,
		RelatedAccountId:     567,
		RelatedAccountAmount: 678,
		HideAmount:           true,
		Comment:              "Scheduled transfer",
		CreatedIp:            "127.0.0.1",
		ScheduledCreated:     true,
	}
	expectedTagIds := []int64{1001, 1002}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, expectedTagIds, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_InvalidTransactionType(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_MODIFY_BALANCE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		TagIds:                 "1001",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.Type = models.TransactionType(255)
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_Weekly(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY,
		ScheduledFrequency:         "0,2,4",
		ScheduledTimezoneUtcOffset: 120,
	}
	location := time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC := time.Date(2026, 11, 30, 22, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	// The local date is Tuesday, while the scheduler's UTC date is Monday.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 120,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	template.ScheduledFrequency = "1,3,5"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	// Sunday is represented by zero.
	template.ScheduledFrequency = "0"
	startTimeInUTC = time.Date(2026, 11, 28, 22, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 11, 28, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_Monthly(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:         "1,15",
		ScheduledTimezoneName:      "Europe/Berlin",
		ScheduledTimezoneUtcOffset: 120,
		CreatedUnixTime:            time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 11, 30, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	// A template created in summer uses the winter offset for December 1.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 60,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	template.ScheduledFrequency = "30"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	// The UTC scheduler stays at 23:00, but the summer transaction is at 22:00.
	template.ScheduledFrequency = "1,15"
	startTimeInUTC = time.Date(2026, 6, 30, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 6, 30, 22, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 120
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_MonthlyLastDays(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "1,-1,-2",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 1, 30, 0, 0, 0, 0, time.UTC)

	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	startTimeInUTC = time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	// February's last day follows the actual length of the month.
	template.ScheduledFrequency = "-1"
	startTimeInUTC = time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	// A positive day that does not exist in February must not create a transaction.
	template.ScheduledFrequency = "31"
	startTimeInUTC = time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_Yearly(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY,
		ScheduledFrequency:         "101,1201",
		ScheduledTimezoneUtcOffset: 120,
	}
	location := time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC := time.Date(2026, 11, 30, 22, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 120,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	template.ScheduledFrequency = "1130"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startTimeInUTC = time.Date(2026, 12, 31, 22, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	template.ScheduledFrequency = "101,1201"
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_YearlyLeapDay(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY,
		ScheduledFrequency:     "229",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)

	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startTimeInUTC = time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_EveryNDaysInvalidFrequency(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "1",
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	// Every N days requires a start date.
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startUnixTime := startTimeInUTC.Unix()
	template.ScheduledStartTime = &startUnixTime
	template.ScheduledFrequency = "1,2"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequency = "0"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledFrequency = "-1"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_EveryNDays(t *testing.T) {
	startUnixTime := time.Date(2026, 11, 28, 0, 0, 0, 0, time.UTC).Unix()
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "3",
		ScheduledStartTime:     &startUnixTime,
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 12, 2, 0, 0, 0, 0, time.UTC)
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	// A negative day difference must be skipped, even if it is divisible by N.
	startTimeInUTC = time.Date(2026, 11, 25, 0, 0, 0, 0, time.UTC)
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	// The start date itself is included.
	startTimeInUTC = time.Date(2026, 11, 28, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_StartAndEndTime(t *testing.T) {
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	startUnixTime := startTimeInUTC.Unix()
	endUnixTime := startTimeInUTC.Unix()
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		ScheduledStartTime:     &startUnixTime,
		ScheduledEndTime:       &endUnixTime,
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)

	// Both boundaries are inclusive.
	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startUnixTime++
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startUnixTime--
	endUnixTime--
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	template.ScheduledStartTime = nil
	endUnixTime++
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledStartTime = &startUnixTime
	template.ScheduledEndTime = nil
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_FixedTimezone(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
	}
	location := time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()
	startTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC)

	// Legacy templates without a timezone name or creation time keep their offset.
	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = 60
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 23, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 60
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = 120
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 22, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 120
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = 180
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 21, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 180
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = 330
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 18, 30, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 330
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = 345
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 18, 15, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 345
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = -300
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 5, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = -300
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = -210
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 3, 30, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = -210
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = -720
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 12, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = -720
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneUtcOffset = 840
	location = time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), location)

	startTimeInUTC = time.Date(2026, 11, 30, 10, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 840
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_InvalidTimezone(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:         "1",
		ScheduledTimezoneName:      "Invalid/Timezone",
		ScheduledTimezoneUtcOffset: 120,
		CreatedUnixTime:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location := time.FixedZone("Template Timezone", int(template.ScheduledTimezoneUtcOffset)*60)
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)

	startTimeInUTC := time.Date(2026, 11, 30, 22, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	// A timezone that cannot be loaded falls back to the stored fixed offset.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 120,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_NamedFixedTimezone(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "1",
		ScheduledTimezoneName:      "UTC",
		ScheduledTimezoneUtcOffset: 480,
		CreatedUnixTime:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := startTimeInUTC.Unix()

	// The named timezone takes precedence over the stored offset.
	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	// Shanghai advances to January 1 across the UTC year boundary.
	template.ScheduledTimezoneName = "Asia/Shanghai"
	location, err = time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC = time.Date(2026, 12, 31, 16, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 480
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneName = "Asia/Kathmandu"
	location, err = time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC = time.Date(2026, 12, 31, 18, 15, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 345
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_TimezoneWithDST(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "1",
		ScheduledTimezoneName:      "Europe/Berlin",
		ScheduledTimezoneUtcOffset: 120,
		CreatedUnixTime:            time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 11, 30, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 60,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 6, 30, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 6, 30, 22, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 120
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_TimezoneWithZeroReferenceOffset(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:         "1",
		ScheduledTimezoneName:      "Europe/London",
		ScheduledTimezoneUtcOffset: 60,
		CreatedUnixTime:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)

	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	// In summer, July 1 midnight is in the previous UTC day.
	startTimeInUTC = time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = time.Date(2026, 6, 30, 23, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 60
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	// Dublin's negative DST uses the same minimum yearly offset.
	template.ScheduledTimezoneName = "Europe/Dublin"
	location, err = time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 0
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_TimezoneWithNegativeOffset(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "1",
		ScheduledTimezoneName:      "America/New_York",
		ScheduledTimezoneUtcOffset: -240,
		CreatedUnixTime:            time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 12, 1, 5, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: -300,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 7, 1, 5, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 7, 1, 4, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = -240
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_SouthernTimezoneWithDST(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "1",
		ScheduledTimezoneName:      "Australia/Sydney",
		ScheduledTimezoneUtcOffset: 660,
		CreatedUnixTime:            time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 6, 30, 14, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 600,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 12, 31, 14, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 12, 31, 13, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 660
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	// Lord Howe advances by thirty minutes rather than one hour.
	template.ScheduledTimezoneName = "Australia/Lord_Howe"
	location, err = time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC = time.Date(2026, 6, 30, 13, 30, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 630
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 12, 31, 13, 30, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 12, 31, 13, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 660
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_TimezoneWithNegativeDST(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                       models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType:     models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:         "1",
		ScheduledTimezoneName:      "Africa/Casablanca",
		ScheduledTimezoneUtcOffset: 60,
		CreatedUnixTime:            time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2019, 5, 10, 0, 0, 0, 0, time.UTC)

	// During Ramadan, Casablanca uses UTC+0 instead of UTC+1.
	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = time.Date(2019, 6, 9, 23, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 60
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_DSTTransition(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "Europe/Berlin",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 3, 28, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 3, 28, 0, 0, 0, 0, time.UTC).Unix()

	// March 29 midnight is before the spring transition.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 60,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 3, 29, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 3, 29, 22, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 120
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	// October 25 midnight is before the autumn transition.
	startTimeInUTC = time.Date(2026, 10, 24, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 10, 24, 22, 0, 0, 0, time.UTC).Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 10, 25, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 60
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_CreatedUnixTime(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "Europe/Moscow",
		CreatedUnixTime:        time.Date(2012, 6, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2014, 11, 30, 20, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2014, 11, 30, 0, 0, 0, 0, time.UTC).Unix()

	// A template created in 2012 retains its UTC+4 reference after the 2014 change.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2014, 11, 30, 21, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: 180,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.CreatedUnixTime = time.Date(2014, 6, 1, 0, 0, 0, 0, time.UTC).Unix()
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)

	startTimeInUTC = time.Date(2014, 11, 30, 21, 0, 0, 0, time.UTC)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	// The creation year is the UTC year, even when Moscow is already in 2014.
	template.CreatedUnixTime = time.Date(2013, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)

	startTimeInUTC = time.Date(2014, 11, 30, 20, 0, 0, 0, time.UTC)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_TimezoneYearBoundary(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "Africa/Sao_Tome",
		CreatedUnixTime:        time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)

	// UTC+0 existed for only the first hour of 2018, but fixes that year's schedule.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2018, 5, 31, 23, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: 60,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	// After the zone returns to UTC+0, the same UTC schedule still has the right date.
	startTimeInUTC = time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 0
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_EveryNDaysTimezoneWithDST(t *testing.T) {
	startUnixTime := time.Date(2026, 3, 27, 23, 0, 0, 0, time.UTC).Unix()
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "2",
		ScheduledStartTime:     &startUnixTime,
		ScheduledTimezoneName:  "Europe/Berlin",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 3, 29, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 3, 29, 0, 0, 0, 0, time.UTC).Unix()

	// March 28 to March 30 is two calendar days, but only 47 elapsed hours.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2026, 3, 29, 22, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: 120,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 3, 30, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC).Unix()
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	// October 24 to October 26 is two calendar days, but 49 elapsed hours.
	startUnixTime = time.Date(2026, 10, 23, 22, 0, 0, 0, time.UTC).Unix()
	startTimeInUTC = time.Date(2026, 10, 25, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 60
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 10, 26, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC).Unix()
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_EveryNDaysDateBoundary(t *testing.T) {
	startUnixTime := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC).Unix()
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "2",
		ScheduledStartTime:     &startUnixTime,
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(time.Now().Unix(), time.UTC)
	startTimeInUTC := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	// February 29 counts as a calendar day in a leap year.
	expectedValue := &models.Transaction{
		Type:             models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:  startTimeInUTC.Unix() * 1000,
		CreatedIp:        "127.0.0.1",
		ScheduledCreated: true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startUnixTime = time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	startTimeInUTC = time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, startTimeInUTC.Unix())
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_StartAndEndTimeTimezoneWithDST(t *testing.T) {
	startUnixTime := time.Date(2026, 6, 30, 22, 0, 0, 0, time.UTC).Unix()
	endUnixTime := startUnixTime
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		ScheduledStartTime:     &startUnixTime,
		ScheduledEndTime:       &endUnixTime,
		ScheduledTimezoneName:  "Europe/Berlin",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 6, 30, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC).Unix()

	// Bounds apply to the final transaction time, not the later scheduler time.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2026, 6, 30, 22, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: 120,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startUnixTime++
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startUnixTime--
	endUnixTime--
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_DateBoundary(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "29",
		ScheduledTimezoneName:  "Asia/Shanghai",
		CreatedUnixTime:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2024, 2, 28, 16, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: 480,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledFrequency = "1"
	startTimeInUTC = time.Date(2024, 2, 29, 16, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 12, 31, 16, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_FutureLeapYear(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "Europe/Berlin",
		CreatedUnixTime:        time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2040, 6, 30, 23, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2040, 6, 30, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2040, 6, 30, 22, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: 120,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	template.ScheduledTimezoneName = "Australia/Sydney"
	location, err = time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC = time.Date(2040, 12, 31, 14, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2040, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2040, 12, 31, 13, 0, 0, 0, time.UTC).Unix() * 1000
	expectedValue.TimezoneUtcOffset = 660
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_ServerTimezone(t *testing.T) {
	originalLocation := time.Local
	t.Cleanup(func() {
		time.Local = originalLocation
	})
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "America/New_York",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 7, 1, 5, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2026, 7, 1, 4, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: -240,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	time.Local = time.UTC
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	time.Local = time.FixedZone("Server UTC-12", -12*60*60)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	time.Local = time.FixedZone("Server UTC-5", -5*60*60)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	time.Local = time.FixedZone("Server UTC+1", 60*60)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	time.Local = time.FixedZone("Server UTC+8", 8*60*60)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	time.Local = time.FixedZone("Server UTC+14", 14*60*60)
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_TimezoneWithMidnightDST(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "America/Santiago",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 9, 5, 4, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: -240,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, _ := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	// September 6 starts at 01:00 UTC-3; midnight must not become September 5.
	startTimeInUTC = time.Date(2026, 9, 6, 4, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = -180
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	startTimeInUTC = time.Date(2026, 9, 7, 4, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = time.Date(2026, 9, 7, 3, 0, 0, 0, time.UTC).Unix() * 1000
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	// Havana's March 8 also starts at 01:00, with UTC-4.
	template.ScheduledTimezoneName = "America/Havana"
	location, err = time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC = time.Date(2026, 3, 8, 5, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = -240
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)

	// Cairo's April 24 starts at 01:00 UTC+3, within the requested local date.
	template.ScheduledTimezoneName = "Africa/Cairo"
	location, err = time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC = time.Date(2026, 4, 23, 22, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 180
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
}

func TestCreateTransactionByScheduledTransactionTemplate_MidnightDSTFrequency(t *testing.T) {
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "6",
		ScheduledTimezoneName:  "America/Santiago",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	location, err := time.LoadLocation(template.ScheduledTimezoneName)
	assert.Equal(t, nil, err)

	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 9, 6, 4, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC).Unix()

	// A missing midnight must not cause the intended day to be skipped.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: -180,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	template.ScheduledFrequencyType = models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY
	template.ScheduledFrequency = "0"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	template.ScheduledFrequencyType = models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY
	template.ScheduledFrequency = "906"
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	startUnixTime := time.Date(2026, 9, 4, 4, 0, 0, 0, time.UTC).Unix()
	template.ScheduledFrequencyType = models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS
	template.ScheduledFrequency = "2"
	template.ScheduledStartTime = &startUnixTime
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_MidnightDSTStartAndEndDate(t *testing.T) {
	location, err := time.LoadLocation("America/Santiago")
	assert.Equal(t, nil, err)
	startTime, err := utils.ParseFromLongDateFirstTimeInTimezone("2026-09-06", location)
	assert.Equal(t, nil, err)
	endTime, err := utils.ParseFromLongDateLastTimeInTimezone("2026-09-06", location)
	assert.Equal(t, nil, err)
	startUnixTime := startTime.Unix()
	endUnixTime := endTime.Unix()
	template := &models.TransactionTemplate{
		TemplateType:           models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS,
		ScheduledFrequency:     "1",
		ScheduledStartTime:     &startUnixTime,
		ScheduledEndTime:       &endUnixTime,
		ScheduledTimezoneName:  "America/Santiago",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 9, 6, 4, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC).Unix()

	// A missing midnight is still the first scheduled calendar date.
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   startTimeInUTC.Unix() * 1000,
		TimezoneUtcOffset: -180,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)

	response := template.ToTransactionTemplateInfoResponse(0)
	assert.Equal(t, "2026-09-06", *response.ScheduledStartDate)
	assert.Equal(t, "2026-09-06", *response.ScheduledEndDate)

	startTimeInUTC = time.Date(2026, 9, 5, 4, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC).Unix()
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startTimeInUTC = time.Date(2026, 9, 7, 4, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC).Unix()
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	// Cairo's start at 01:00 must not move the returned end date to April 25.
	location, err = time.LoadLocation("Africa/Cairo")
	assert.Equal(t, nil, err)
	startTime, err = utils.ParseFromLongDateFirstTimeInTimezone("2026-04-24", location)
	assert.Equal(t, nil, err)
	endTime, err = utils.ParseFromLongDateLastTimeInTimezone("2026-04-24", location)
	assert.Equal(t, nil, err)
	startUnixTime = startTime.Unix()
	endUnixTime = endTime.Unix()
	template.ScheduledTimezoneName = "Africa/Cairo"
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC = time.Date(2026, 4, 23, 22, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue.TransactionTime = startTimeInUTC.Unix() * 1000
	expectedValue.TimezoneUtcOffset = 180
	actualValue, _ = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	response = template.ToTransactionTemplateInfoResponse(0)
	assert.Equal(t, "2026-04-24", *response.ScheduledStartDate)
	assert.Equal(t, "2026-04-24", *response.ScheduledEndDate)
}

func TestCreateTransactionByScheduledTransactionTemplate_AmbiguousMidnight(t *testing.T) {
	location, err := time.LoadLocation("America/Havana")
	assert.Equal(t, nil, err)
	startTime, err := utils.ParseFromLongDateFirstTimeInTimezone("2026-11-01", location)
	assert.Equal(t, nil, err)
	startUnixTime := startTime.Unix()
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY,
		ScheduledFrequency:     "1",
		ScheduledStartTime:     &startUnixTime,
		ScheduledTimezoneName:  "America/Havana",
		CreatedUnixTime:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2026, 11, 1, 5, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC).Unix()

	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2026, 11, 1, 4, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: -240,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)
}

func TestCreateTransactionByScheduledTransactionTemplate_SkippedDate(t *testing.T) {
	location, err := time.LoadLocation("Pacific/Apia")
	assert.Equal(t, nil, err)
	template := &models.TransactionTemplate{
		Type:                   models.TRANSACTION_TYPE_EXPENSE,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY,
		ScheduledFrequency:     "1",
		ScheduledTimezoneName:  "Pacific/Apia",
		CreatedUnixTime:        time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	template.ScheduledAt = utils.GetScheduledTransactionScheduledAt(template.CreatedUnixTime, location)
	startTimeInUTC := time.Date(2011, 12, 30, 10, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC := time.Date(2011, 12, 30, 0, 0, 0, 0, time.UTC).Unix()

	// A skipped date must not create an extra transaction for December 31.
	actualValue, actualTagIds := Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Nil(t, actualValue)
	assert.Nil(t, actualTagIds)

	startTimeInUTC = time.Date(2011, 12, 31, 10, 0, 0, 0, time.UTC)
	todayFirstUnixTimeInUTC = time.Date(2011, 12, 31, 0, 0, 0, 0, time.UTC).Unix()
	expectedValue := &models.Transaction{
		Type:              models.TRANSACTION_DB_TYPE_EXPENSE,
		TransactionTime:   time.Date(2011, 12, 30, 10, 0, 0, 0, time.UTC).Unix() * 1000,
		TimezoneUtcOffset: 840,
		CreatedIp:         "127.0.0.1",
		ScheduledCreated:  true,
	}
	actualValue, actualTagIds = Transactions.createTransactionByScheduledTransactionTemplate(core.NewNullContext(), template, startTimeInUTC, todayFirstUnixTimeInUTC)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, []int64{}, actualTagIds)
}
