package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/mcp"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

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
	service := &TokensApi{}

	expectedValue := int64(80)
	actualValue, err := service.getMCPToolsMask(core.NewNullContext(), []string{"add_transaction", "query_latest_exchange_rates"})

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTokensApiGetMCPToolsMask_EmptyNames(t *testing.T) {
	initializeTokensApiTest(t)
	service := &TokensApi{}

	expectedValue := int64(0)
	actualValue, err := service.getMCPToolsMask(core.NewNullContext(), []string{})

	assert.Equal(t, errs.ErrMCPToolsEmpty, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTokensApiGetMCPToolsMask_DuplicateNames(t *testing.T) {
	initializeTokensApiTest(t)
	service := &TokensApi{}

	expectedValue := int64(32)
	actualValue, err := service.getMCPToolsMask(core.NewNullContext(), []string{"query_transactions", "query_transactions"})

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedValue, actualValue)
}

func TestTokensApiGetMCPToolsMask_InvalidNames(t *testing.T) {
	initializeTokensApiTest(t)
	service := &TokensApi{}

	actualValue, err := service.getMCPToolsMask(core.NewNullContext(), []string{"unknown_tool"})

	assert.Equal(t, int64(0), actualValue)
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, err)

	actualValue, err = service.getMCPToolsMask(core.NewNullContext(), []string{"add_transaction", ""})

	assert.Equal(t, int64(0), actualValue)
	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, err)
}
