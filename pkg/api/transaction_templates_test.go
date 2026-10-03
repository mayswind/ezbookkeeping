package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
	"github.com/mayswind/ezbookkeeping/pkg/validators"
)

func newScheduledTemplateTestDatabase(t *testing.T) (core.Context, *datastore.Database) {
	t.Helper()
	c := core.NewNullContext()
	previousDataStore := *datastore.Container
	previousUuid := *uuid.Container
	config := &settings.Config{
		DatabaseConfig: &settings.DatabaseConfig{
			DatabaseType:      settings.Sqlite3DbType,
			DatabasePath:      filepath.Join(t.TempDir(), "scheduled.db"),
			MaxOpenConnection: 1,
		},
		UuidGeneratorType: settings.InternalUuidGeneratorType,
	}
	require.NoError(t, datastore.InitializeDataStore(config))
	require.NoError(t, uuid.InitializeUuidGenerator(config))
	t.Cleanup(func() {
		sess := datastore.Container.UserDataStore.Choose(1).NewSession(c)
		engine := sess.Engine()
		require.NoError(t, sess.Close())
		require.NoError(t, engine.Close())
		*datastore.Container = previousDataStore
		*uuid.Container = previousUuid
	})
	require.NoError(t, datastore.Container.UserDataStore.SyncStructs(
		&models.User{}, &models.TransactionTemplate{}, &models.Transaction{}, &models.Account{},
		&models.TransactionCategory{}, &models.TransactionTag{}, &models.TransactionTagIndex{}, &models.TransactionPictureInfo{},
	))
	db := datastore.Container.UserDataStore.Choose(1)
	sess := db.NewSession(c)
	_, err := sess.Insert(&models.Account{
		AccountId: 1, Uid: 1, Type: models.ACCOUNT_TYPE_SINGLE_ACCOUNT, Currency: "EUR",
	}, &models.TransactionCategory{
		CategoryId: 10, Uid: 1, Type: models.CATEGORY_TYPE_EXPENSE,
	}, &models.TransactionCategory{
		CategoryId: 11, Uid: 1, ParentCategoryId: 10, Type: models.CATEGORY_TYPE_EXPENSE,
	})
	require.NoError(t, err)
	require.NoError(t, sess.Close())

	return c, db
}

func TestScheduledTemplateAfterDaylightSavingTime(t *testing.T) {
	c, db := newScheduledTemplateTestDatabase(t)
	validate := binding.Validator.Engine().(*validator.Validate)
	require.NoError(t, validate.RegisterValidation("validTagFilter", validators.ValidTagFilter))
	require.NoError(t, validate.RegisterValidation("validAmountFilter", validators.ValidAmountFilter))
	userSession := db.NewSession(c)
	_, err := userSession.Insert(&models.User{Uid: 1, Username: "user", Email: "user@example.com"})
	require.NoError(t, err)
	require.NoError(t, userSession.Close())

	var request models.TransactionTemplateCreateRequest
	require.NoError(t, json.Unmarshal([]byte(`{
		"templateType": 2, "name": "Monthly bill", "type": 3,
		"categoryId": "11", "sourceAccountId": "1", "sourceAmount": 1250,
		"scheduledFrequencyType": 2, "scheduledFrequency": "1",
		"utcOffset": 120, "timeZone": "Europe/Berlin"
	}`), &request))
	template, err := TransactionTemplates.createNewTemplateModel(1, &request, 0)
	require.NoError(t, err)
	assert.Equal(t, int16(1320), template.ScheduledAt)
	template.TemplateId = 1
	sess := db.NewSession(c)
	_, err = sess.Insert(template)
	require.NoError(t, err)
	require.NoError(t, sess.Close())

	// The template was saved at UTC+2, but Berlin is at UTC+1 in December.
	for _, step := range []struct {
		utcHour int
		count   int
	}{
		{22, 0},
		{23, 1},
	} {
		runTime := time.Date(2026, 11, 30, step.utcHour, 0, 0, 0, time.UTC)
		require.NoError(t, services.Transactions.CreateScheduledTransactions(c, runTime.Unix(), 15*time.Minute))
		var transactions []*models.Transaction
		sess = db.NewSession(c)
		require.NoError(t, sess.Find(&transactions))
		require.NoError(t, sess.Close())
		assert.Len(t, transactions, step.count, "run at %s", runTime)
		for _, transaction := range transactions {
			unixTime := utils.GetUnixTimeFromTransactionTime(transaction.TransactionTime)
			t.Logf("run=%s transaction=%s offset=%d", runTime, time.Unix(unixTime, 0).UTC(), transaction.TimezoneUtcOffset)
			assert.Equal(t, time.Date(2026, 11, 30, 23, 0, 0, 0, time.UTC).Unix(), unixTime)
			assert.Equal(t, int16(60), transaction.TimezoneUtcOffset)
			assert.True(t, transaction.ScheduledCreated)
		}
	}

	// These are the two list requests in the reproduction posted in #566.
	for _, query := range []struct {
		date  string
		count int
	}{
		{"2026-12-01", 1},
		{"2026-11-30", 0},
	} {
		start, err := utils.ParseFromLongDateFirstTime(query.date, 60)
		require.NoError(t, err)
		end, err := utils.ParseFromLongDateLastTime(query.date, 60)
		require.NoError(t, err)
		url := fmt.Sprintf("/transactions/list.json?count=10&page=1&min_time=%d&max_time=%d", utils.GetMinTransactionTimeFromUnixTime(start.Unix()), utils.GetMaxTransactionTimeFromUnixTime(end.Unix()))
		ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		ginContext.Request = httptest.NewRequest("GET", url, nil)
		ginContext.Request.Header.Set("X-Timezone-Offset", "60")
		webContext := core.WrapWebContext(ginContext, nil)
		webContext.SetTokenClaims(&core.UserTokenClaims{Uid: 1})
		response, apiErr := Transactions.TransactionListHandler(webContext)
		require.Nil(t, apiErr)
		items := response.(*models.TransactionInfoPageWrapperResponse).Items
		t.Logf("transactions/list.json %s X-Timezone-Offset: 60 -> %d items", query.date, len(items))
		assert.Len(t, items, query.count, "list date %s", query.date)
	}
}

func TestScheduledTemplateTimezoneBoundaries(t *testing.T) {
	for _, test := range []struct {
		name      string
		zone      string
		offset    int16
		frequency string
		kind      models.TransactionScheduleFrequencyType
		start     string
		end       string
		run       string
		count     int
		actual    string
		actualUTC int16
		nullZone  bool
	}{
		{name: "summer", zone: "Europe/Berlin", offset: 60, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, frequency: "1", run: "2026-06-30T22:00:00Z", count: 1, actual: "2026-06-30T22:00:00Z", actualUTC: 120},
		{name: "legacy empty", offset: 120, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, frequency: "1", run: "2026-11-30T22:00:00Z", count: 1, actual: "2026-11-30T22:00:00Z", actualUTC: 120},
		{name: "legacy null", offset: 120, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, frequency: "1", run: "2026-11-30T22:00:00Z", count: 1, actual: "2026-11-30T22:00:00Z", actualUTC: 120, nullZone: true},
		{name: "start date", zone: "Europe/Berlin", offset: 120, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, frequency: "1", start: "2026-12-01", end: "2026-12-01", run: "2026-11-30T23:00:00Z", count: 1, actual: "2026-11-30T23:00:00Z", actualUTC: 60},
		{name: "before start", zone: "Europe/Berlin", offset: 120, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, frequency: "1", start: "2026-12-02", run: "2026-11-30T23:00:00Z"},
		{name: "after end", zone: "Europe/Berlin", offset: 120, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, frequency: "1", end: "2026-11-30", run: "2026-11-30T23:00:00Z"},
		{name: "calendar days across spring DST", zone: "Europe/Berlin", offset: 60, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, frequency: "2", start: "2026-03-28", run: "2026-03-29T22:00:00Z", count: 1, actual: "2026-03-29T22:00:00Z", actualUTC: 120},
		{name: "off day across spring DST", zone: "Europe/Berlin", offset: 60, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_EVERY_N_DAYS, frequency: "2", start: "2026-03-28", run: "2026-03-30T22:00:00Z"},
		{name: "quarter hour zone", zone: "Asia/Kathmandu", offset: 345, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, frequency: "1", run: "2026-11-30T18:15:00Z", count: 1, actual: "2026-11-30T18:15:00Z", actualUTC: 345},
		{name: "exclusive interval end", zone: "Asia/Kathmandu", offset: 345, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, frequency: "1", run: "2026-11-30T18:00:00Z"},
		{name: "missing midnight", zone: "America/Sao_Paulo", offset: -180, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, frequency: "1", run: "2018-11-04T03:00:00Z", count: 1, actual: "2018-11-04T03:00:00Z", actualUTC: -120},
		{name: "repeated midnight first occurrence", zone: "America/Havana", offset: -240, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, frequency: "1", run: "2026-11-01T04:00:00Z", count: 1, actual: "2026-11-01T04:00:00Z", actualUTC: -240},
		{name: "repeated midnight second occurrence", zone: "America/Havana", offset: -240, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY, frequency: "1", run: "2026-11-01T05:00:00Z"},
		{name: "weekly local day", zone: "Europe/Berlin", offset: 120, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_WEEKLY, frequency: "2", run: "2026-11-30T23:00:00Z", count: 1, actual: "2026-11-30T23:00:00Z", actualUTC: 60},
		{name: "yearly local day", zone: "Europe/Berlin", offset: 120, kind: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_YEARLY, frequency: "1201", run: "2026-11-30T23:00:00Z", count: 1, actual: "2026-11-30T23:00:00Z", actualUTC: 60},
	} {
		t.Run(test.name, func(t *testing.T) {
			c, db := newScheduledTemplateTestDatabase(t)
			request := models.TransactionTemplateCreateRequest{
				TemplateType: models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE, Name: "Bill", Type: models.TRANSACTION_TYPE_EXPENSE,
				CategoryId: 11, SourceAccountId: 1, SourceAmount: 1250,
				ScheduledFrequencyType: &test.kind, ScheduledFrequency: &test.frequency,
				ScheduledTimezoneUtcOffset: &test.offset, ScheduledTimezoneName: &test.zone,
			}
			if test.start != "" {
				request.ScheduledStartDate = &test.start
			}
			if test.end != "" {
				request.ScheduledEndDate = &test.end
			}
			template, err := TransactionTemplates.createNewTemplateModel(1, &request, 0)
			require.NoError(t, err)
			template.TemplateId = 1
			sess := db.NewSession(c)
			_, err = sess.Insert(template)
			require.NoError(t, err)
			if test.nullZone {
				_, err = sess.Exec("UPDATE transaction_template SET scheduled_timezone_name=NULL WHERE template_id=1")
				require.NoError(t, err)
			}
			require.NoError(t, sess.Close())
			run, err := time.Parse(time.RFC3339, test.run)
			require.NoError(t, err)
			require.NoError(t, services.Transactions.CreateScheduledTransactions(c, run.Unix(), 15*time.Minute))
			var transactions []*models.Transaction
			sess = db.NewSession(c)
			require.NoError(t, sess.Find(&transactions))
			require.NoError(t, sess.Close())
			require.Len(t, transactions, test.count)
			if test.count > 0 {
				actual, err := time.Parse(time.RFC3339, test.actual)
				require.NoError(t, err)
				assert.Equal(t, actual.Unix(), utils.GetUnixTimeFromTransactionTime(transactions[0].TransactionTime))
				assert.Equal(t, test.actualUTC, transactions[0].TimezoneUtcOffset)
				assert.True(t, transactions[0].ScheduledCreated)
			}
		})
	}
}

func TestScheduledTemplateTimezoneValidation(t *testing.T) {
	for _, zone := range []string{"Europe/Berlin", "", "Local", "No/SuchZone"} {
		t.Run(zone, func(t *testing.T) {
			frequencyType := models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_DAILY
			frequency := "1"
			offset := int16(120)
			request := models.TransactionTemplateCreateRequest{
				TemplateType:           models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
				ScheduledFrequencyType: &frequencyType, ScheduledFrequency: &frequency,
				ScheduledTimezoneUtcOffset: &offset, ScheduledTimezoneName: &zone,
			}
			template, err := TransactionTemplates.createNewTemplateModel(1, &request, 0)
			if zone == "Local" || zone == "No/SuchZone" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			response := template.ToTransactionTemplateInfoResponse(0)
			if zone == "" {
				assert.Nil(t, response.ScheduledTimezoneName)
			} else {
				require.NotNil(t, response.ScheduledTimezoneName)
				assert.Equal(t, zone, *response.ScheduledTimezoneName)
			}
		})
	}
}

func TestScheduledTemplateModifyTimezone(t *testing.T) {
	c, db := newScheduledTemplateTestDatabase(t)
	previousConfig := settings.Container.GetCurrentConfig()
	settings.SetCurrentConfig(&settings.Config{EnableScheduledTransaction: true})
	t.Cleanup(func() { settings.SetCurrentConfig(previousConfig) })
	validate := binding.Validator.Engine().(*validator.Validate)
	require.NoError(t, validate.RegisterValidation("notBlank", validators.NotBlank))
	require.NoError(t, validate.RegisterValidation("validTransactionAmount", validators.ValidTransactionAmount))
	sess := db.NewSession(c)
	_, err := sess.Insert(&models.TransactionTemplate{
		TemplateId: 1, Uid: 1, TemplateType: models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		Name: "Bill", Type: models.TRANSACTION_TYPE_EXPENSE, AccountId: 1, CategoryId: 11, Amount: 1250,
		ScheduledFrequencyType: models.TRANSACTION_SCHEDULE_FREQUENCY_TYPE_MONTHLY, ScheduledFrequency: "1",
		ScheduledAt: 1320, ScheduledTimezoneUtcOffset: 120, ScheduledTimezoneName: "Europe/Berlin",
	})
	require.NoError(t, err)
	require.NoError(t, sess.Close())
	for index, test := range []struct {
		zone *string
		want string
	}{
		{nil, "Europe/Berlin"},
		{new("Asia/Kathmandu"), "Asia/Kathmandu"},
		{new(""), ""},
		{nil, ""},
	} {
		request := map[string]any{
			"id": "1", "name": fmt.Sprintf("Bill %d", index), "type": 3,
			"categoryId": "11", "sourceAccountId": "1", "sourceAmount": 1250,
			"scheduledFrequencyType": 2, "scheduledFrequency": "1", "utcOffset": 120,
		}
		if test.zone != nil {
			request["timeZone"] = *test.zone
		}
		payload, err := json.Marshal(request)
		require.NoError(t, err)
		ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		ginContext.Request = httptest.NewRequest("POST", "/", bytes.NewReader(payload))
		ginContext.Request.Header.Set("Content-Type", "application/json")
		ginContext.Request.Header.Set("X-Timezone-Name", "America/New_York")
		webContext := core.WrapWebContext(ginContext, nil)
		webContext.SetTokenClaims(&core.UserTokenClaims{Uid: 1})
		response, apiErr := TransactionTemplates.TemplateModifyHandler(webContext)
		require.Nil(t, apiErr)
		info := response.(*models.TransactionTemplateInfoResponse)
		if test.want == "" {
			assert.Nil(t, info.ScheduledTimezoneName)
		} else {
			require.NotNil(t, info.ScheduledTimezoneName)
			assert.Equal(t, test.want, *info.ScheduledTimezoneName)
		}
		stored := &models.TransactionTemplate{}
		sess = db.NewSession(c)
		found, err := sess.ID(1).Get(stored)
		require.NoError(t, err)
		require.True(t, found)
		require.NoError(t, sess.Close())
		assert.Equal(t, test.want, stored.ScheduledTimezoneName)
	}
}

func TestScheduledTemplateDateAcrossDaylightSavingTime(t *testing.T) {
	template := &models.TransactionTemplate{
		TemplateType:          models.TRANSACTION_TEMPLATE_TYPE_SCHEDULE,
		ScheduledTimezoneName: "Europe/Berlin", ScheduledTimezoneUtcOffset: 120,
	}
	for _, test := range []struct {
		date  string
		start string
		end   string
	}{
		{"2026-03-29", "2026-03-28T23:00:00Z", "2026-03-29T21:59:59Z"},
		{"2026-10-25", "2026-10-24T22:00:00Z", "2026-10-25T22:59:59Z"},
	} {
		t.Run(test.date, func(t *testing.T) {
			start, err := template.ParseScheduledDate(test.date, false)
			require.NoError(t, err)
			end, err := template.ParseScheduledDate(test.date, true)
			require.NoError(t, err)
			assert.Equal(t, test.start, start.UTC().Format(time.RFC3339))
			assert.Equal(t, test.end, end.UTC().Truncate(time.Second).Format(time.RFC3339))
			startUnix, endUnix := start.Unix(), end.Unix()
			template.ScheduledStartTime, template.ScheduledEndTime = &startUnix, &endUnix
			response := template.ToTransactionTemplateInfoResponse(0)
			assert.Equal(t, test.date, *response.ScheduledStartDate)
			assert.Equal(t, test.date, *response.ScheduledEndDate)
		})
	}
}
