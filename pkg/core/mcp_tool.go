package core

// MCPToolType represents the type of MCP tool
type MCPToolType int64

const (
	MCP_TOOL_TYPE_QUERY_ALL_ACCOUNTS               MCPToolType = 1
	MCP_TOOL_TYPE_QUERY_ALL_ACCOUNTS_BALANCE       MCPToolType = 2
	MCP_TOOL_TYPE_QUERY_ALL_TRANSACTION_CATEGORIES MCPToolType = 3
	MCP_TOOL_TYPE_QUERY_ALL_TRANSACTION_TAGS       MCPToolType = 4
	MCP_TOOL_TYPE_QUERY_LATEST_EXCHANGE_RATES      MCPToolType = 5
	MCP_TOOL_TYPE_QUERY_TRANSACTIONS               MCPToolType = 6
	MCP_TOOL_TYPE_ADD_TRANSACTION                  MCPToolType = 7
	MCP_TOOL_TYPE_DELETE_TRANSACTION               MCPToolType = 8

	MCP_TOOL_MAX_TYPE_VALUE = MCP_TOOL_TYPE_DELETE_TRANSACTION

	// MCPAllTools is the mask of all currently supported tools
	MCPAllTools int64 = (1 << MCP_TOOL_MAX_TYPE_VALUE) - 1
)

// IsEnabled returns whether the tool is enabled in the given permission mask
func (t MCPToolType) IsEnabled(mask int64) bool {
	return t >= 1 && mask >= 0 && mask&(int64(1)<<(t-1)) != 0
}
