package mcp

import (
	stdjson "encoding/json"
	"encoding/json/jsontext"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestMCPDeleteTransactionHandle_NilArguments(t *testing.T) {
	_, _, err := MCPDeleteTransactionToolHandler.Handle(nil, &MCPCallToolRequest{
		Name: "delete_transaction",
	}, nil, nil, nil)

	assert.Equal(t, errs.ErrIncompleteOrIncorrectSubmission, err)
}

func TestMCPDeleteTransactionHandle_MalformedArguments(t *testing.T) {
	_, _, err := MCPDeleteTransactionToolHandler.Handle(nil, &MCPCallToolRequest{
		Name:      "delete_transaction",
		Arguments: jsontext.Value(`{"transaction_id":`),
	}, nil, nil, nil)

	assert.EqualError(t, err, errs.ErrIncompleteOrIncorrectSubmission.Message)
}

func TestMCPDeleteTransactionHandle_InvalidTransactionId(t *testing.T) {
	testCases := []string{
		`{"transaction_id":"not-a-number"}`,
		`{"transaction_id":"0"}`,
	}

	for _, testCase := range testCases {
		_, _, err := MCPDeleteTransactionToolHandler.Handle(nil, &MCPCallToolRequest{
			Name:      "delete_transaction",
			Arguments: jsontext.Value(testCase),
		}, nil, nil, nil)

		assert.Equal(t, errs.ErrTransactionIdInvalid, err, testCase)
	}
}

func TestMCPDeleteTransactionResponse_TransactionIdIsString(t *testing.T) {
	structuredResponse, content, err := MCPDeleteTransactionToolHandler.createMCPDeleteTransactionResponse(3849953921556545536, true)

	assert.NoError(t, err)

	response, ok := structuredResponse.(MCPDeleteTransactionResponse)
	assert.True(t, ok)
	assert.True(t, response.Success)
	assert.True(t, response.DryRun)
	assert.Equal(t, "3849953921556545536", response.TransactionId)

	if assert.Len(t, content, 1) {
		var textResponse map[string]any
		assert.NoError(t, stdjson.Unmarshal([]byte(content[0].Text), &textResponse))
		assert.Equal(t, "3849953921556545536", textResponse["transaction_id"])
		assert.Equal(t, true, textResponse["dry_run"])
	}
}

func TestMCPDeleteTransactionResponse_DryRunIsOmittedWhenFalse(t *testing.T) {
	_, content, err := MCPDeleteTransactionToolHandler.createMCPDeleteTransactionResponse(1, false)

	assert.NoError(t, err)

	if assert.Len(t, content, 1) {
		var textResponse map[string]any
		assert.NoError(t, stdjson.Unmarshal([]byte(content[0].Text), &textResponse))
		assert.Equal(t, "1", textResponse["transaction_id"])
		assert.NotContains(t, textResponse, "dry_run")
	}
}
