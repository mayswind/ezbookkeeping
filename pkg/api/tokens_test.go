package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/mcp"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

type tokenApiTestData struct {
	apiTestModelOperations
	records []*models.TokenRecord
}

func (d *tokenApiTestData) Get(s *apiTestSession, beans ...any) (bool, error) {
	if len(beans) != 1 {
		return false, s.unsupportedOperation("Get", beans)
	}

	switch row := beans[0].(type) {
	case *models.TokenRecord:
		if row == nil || s.query != "uid=? AND user_token_id=? AND created_unix_time=?" || len(s.args) != 3 {
			return false, s.unsupportedOperation("Get", row)
		}

		uid, uidValid := s.args[0].(int64)
		userTokenId, userTokenIdValid := s.args[1].(int64)
		createdUnixTime, createdUnixTimeValid := s.args[2].(int64)

		if !uidValid || !userTokenIdValid || !createdUnixTimeValid {
			return false, s.unsupportedOperation("Get", row)
		}

		for _, record := range d.records {
			if record != nil && record.Uid == uid && record.UserTokenId == userTokenId && record.CreatedUnixTime == createdUnixTime {
				*row = *record
				return true, nil
			}
		}

		return false, nil
	default:
		return false, s.unsupportedOperation("Get", beans[0])
	}
}

func (d *tokenApiTestData) Insert(s *apiTestSession, beans ...any) (int64, error) {
	if len(beans) != 1 || s.query != nil || len(s.args) != 0 {
		return 0, s.unsupportedOperation("Insert", beans)
	}

	switch record := beans[0].(type) {
	case *models.TokenRecord:
		if record == nil {
			return 0, s.unsupportedOperation("Insert", record)
		}

		recordCopy := *record
		d.records = append(d.records, &recordCopy)
		return 1, nil
	default:
		return 0, s.unsupportedOperation("Insert", beans[0])
	}
}

func initializeTokensApiTest(t *testing.T) {
	t.Helper()
	previousContainer := mcp.Container
	t.Cleanup(func() {
		mcp.Container = previousContainer
	})

	err := mcp.InitializeMCPHandlers(&settings.Config{})

	if err != nil {
		t.Fatal(err)
	}
}

func TestTokensApiGetMCPToolsMask(t *testing.T) {
	initializeTokensApiTest(t)
	api := &TokensApi{}

	expectedValue := int64(80)
	actualValue, err := api.getMCPToolsMask(core.NewNullContext(), []string{"add_transaction", "query_latest_exchange_rates"})

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTokensApiGetMCPToolsMask_EmptyNames(t *testing.T) {
	initializeTokensApiTest(t)
	api := &TokensApi{}

	expectedValue := int64(0)
	actualValue, err := api.getMCPToolsMask(core.NewNullContext(), []string{})

	assert.Equal(t, errs.ErrMCPToolsEmpty, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTokensApiGetMCPToolsMask_DuplicateNames(t *testing.T) {
	initializeTokensApiTest(t)
	api := &TokensApi{}

	expectedValue := int64(32)
	actualValue, err := api.getMCPToolsMask(core.NewNullContext(), []string{"query_transactions", "query_transactions"})

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTokensApiGetMCPToolsMask_InvalidNames(t *testing.T) {
	initializeTokensApiTest(t)
	api := &TokensApi{}

	actualValue, err := api.getMCPToolsMask(core.NewNullContext(), []string{"unknown_tool"})

	assert.Equal(t, int64(0), actualValue)
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, err)

	actualValue, err = api.getMCPToolsMask(core.NewNullContext(), []string{"add_transaction", ""})

	assert.Equal(t, int64(0), actualValue)
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, err)
}
