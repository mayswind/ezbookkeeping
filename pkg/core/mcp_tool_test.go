package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMCPToolTypeIsEnabled(t *testing.T) {
	assert.True(t, MCP_TOOL_TYPE_ADD_TRANSACTION.IsEnabled(MCPAllTools))
	assert.True(t, MCP_TOOL_TYPE_QUERY_LATEST_EXCHANGE_RATES.IsEnabled(MCPAllTools))
	assert.False(t, MCP_TOOL_TYPE_ADD_TRANSACTION.IsEnabled(0))

	assert.True(t, MCP_TOOL_TYPE_ADD_TRANSACTION.IsEnabled(65))
	assert.False(t, MCP_TOOL_TYPE_QUERY_TRANSACTIONS.IsEnabled(65))
	assert.True(t, MCP_TOOL_TYPE_QUERY_LATEST_EXCHANGE_RATES.IsEnabled(17))
	assert.True(t, MCP_TOOL_TYPE_QUERY_ALL_TRANSACTION_TAGS.IsEnabled(127))
}

func TestMCPToolTypeIsEnabled_InvalidMask(t *testing.T) {
	assert.False(t, MCP_TOOL_TYPE_ADD_TRANSACTION.IsEnabled(-1))
}

func TestMCPToolTypeIsEnabled_InvalidToolType(t *testing.T) {
	assert.False(t, MCPToolType(0).IsEnabled(MCPAllTools))
	assert.False(t, MCPToolType(4503599627370496).IsEnabled(MCPAllTools))
}

func TestMCPToolTypeIsEnabled_HighBit(t *testing.T) {
	assert.True(t, MCPToolType(63).IsEnabled(4611686018427387904))
}
