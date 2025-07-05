# MCP (Model Context Protocol) Support for ezBookkeeping

This document describes the MCP (Model Context Protocol) implementation in ezBookkeeping, which allows AI models and LLMs to securely access and query transaction data.

## What is MCP?

Model Context Protocol (MCP) is an open standard developed by Anthropic that enables AI models to securely connect to external data sources and tools. It provides a standardized way for AI assistants to access information and perform actions while maintaining security and user control.

## Features

ezBookkeeping's MCP server provides the following capabilities:

### Resources
- **Recent Transactions**: Access to the latest transactions
- **Account List**: Complete list of user accounts with balances
- **Category List**: All transaction categories
- **Tag List**: All transaction tags

### Tools
- **query_transactions**: Advanced transaction querying with filters
  - Filter by date range, category, account, transaction type
  - Search by keywords
  - Limit results
- **get_transaction_statistics**: Get transaction counts and summaries
  - Breakdown by transaction type (income, expense, transfer)
  - Filtered by date range, category, or account
- **get_account_balance**: Retrieve account balances
  - Individual account or all accounts
  - Current balance information

## Configuration

### Enable MCP Server

Add the following to your `ezbookkeeping.ini` configuration file:

```ini
[mcp]
# Set to true to enable MCP (Model Context Protocol) server for AI/LLM access
enable = true
```

### Environment Variable

You can also enable MCP using an environment variable:

```bash
export EBK_MCP_ENABLE=true
```

## API Endpoint

When MCP is enabled, the server exposes an MCP endpoint at:

```
POST /api/v1/mcp
```

This endpoint implements the JSON-RPC 2.0 protocol required by MCP.

## Authentication

MCP requests require the same authentication as other API endpoints:
- Valid JWT token in Authorization header
- Or valid session cookie

## Usage Examples

### Initialize MCP Connection

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "resources": {},
      "tools": {}
    },
    "clientInfo": {
      "name": "my-ai-assistant",
      "version": "1.0.0"
    }
  }
}
```

### List Available Resources

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "resources/list"
}
```

### Read Recent Transactions

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "resources/read",
  "params": {
    "uri": "transactions://recent"
  }
}
```

### Query Transactions with Filters

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "query_transactions",
    "arguments": {
      "start_date": "2024-01-01",
      "end_date": "2024-12-31",
      "type": "expense",
      "limit": 50
    }
  }
}
```

### Get Transaction Statistics

```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "get_transaction_statistics",
    "arguments": {
      "start_date": "2024-01-01",
      "end_date": "2024-12-31"
    }
  }
}
```

### Get Account Balances

```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "get_account_balance"
  }
}
```

## Security Considerations

- MCP access requires the same authentication as regular API access
- Users must explicitly enable MCP in their configuration
- All data access is scoped to the authenticated user only
- No write operations are exposed through MCP (read-only access)
- Error messages do not expose sensitive information

## Integration with AI Tools

### Claude Desktop

To use with Claude Desktop, add to your configuration:

```json
{
  "mcpServers": {
    "ezbookkeeping": {
      "command": "curl",
      "args": [
        "-X", "POST",
        "-H", "Content-Type: application/json",
        "-H", "Authorization: Bearer YOUR_TOKEN",
        "http://localhost:8080/api/v1/mcp"
      ]
    }
  }
}
```

### Other MCP Clients

Any MCP-compatible client can connect to the ezBookkeeping MCP server using the standard MCP protocol over HTTP.

## Troubleshooting

### MCP Not Available

If you get "MCP is not enabled" error:
1. Check that `enable = true` is set in the `[mcp]` section of your config
2. Restart the ezBookkeeping server
3. Verify you're using a valid authentication token

### Authentication Errors

Ensure you're including valid authentication:
- JWT token in `Authorization: Bearer <token>` header
- Or valid session cookies

### Connection Issues

Check that:
- ezBookkeeping server is running
- Port 8080 (or your configured port) is accessible
- No firewall is blocking the connection

## Development

### Adding New Resources

To add a new resource:

1. Add the resource definition in `handleListResources`
2. Add the resource handler in `handleReadResource`
3. Implement the data retrieval method

### Adding New Tools

To add a new tool:

1. Add the tool definition in `handleListTools`
2. Add the tool handler in `handleCallTool`
3. Implement the tool execution method

### Testing

Run the MCP tests:

```bash
go test -v ./pkg/mcp/
```

## Protocol Compliance

This implementation follows the MCP protocol specification:
- JSON-RPC 2.0 transport
- Protocol version 2024-11-05
- Standard resource and tool capabilities
- Proper error handling and responses