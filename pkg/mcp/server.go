package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const (
	MCPVersion = "2024-11-05"
)

// MCPServer represents the MCP server
type MCPServer struct {
	transactions          *services.TransactionService
	accounts              *services.AccountService
	transactionCategories *services.TransactionCategoryService
	transactionTags       *services.TransactionTagService
	users                 *services.UserService
	config                *settings.Config
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	return &MCPServer{
		transactions:          services.Transactions,
		accounts:              services.Accounts,
		transactionCategories: services.TransactionCategories,
		transactionTags:       services.TransactionTags,
		users:                 services.Users,
		config:                settings.Container.Current,
	}
}

// JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP protocol structures
type MCPServerInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type MCPClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type MCPCapabilities struct {
	Resources *MCPResourceCapabilities `json:"resources,omitempty"`
	Tools     *MCPToolCapabilities     `json:"tools,omitempty"`
	Prompts   *MCPPromptCapabilities   `json:"prompts,omitempty"`
}

type MCPResourceCapabilities struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type MCPToolCapabilities struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type MCPPromptCapabilities struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type MCPInitializeRequest struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    MCPCapabilities `json:"capabilities"`
	ClientInfo      MCPClientInfo   `json:"clientInfo"`
}

type MCPInitializeResponse struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    MCPCapabilities `json:"capabilities"`
	ServerInfo      MCPServerInfo   `json:"serverInfo"`
}

// Resource structures
type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

type MCPResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"`
}

type MCPListResourcesResponse struct {
	Resources []MCPResource `json:"resources"`
	NextCursor string       `json:"nextCursor,omitempty"`
}

type MCPReadResourceRequest struct {
	URI string `json:"uri"`
}

type MCPReadResourceResponse struct {
	Contents []MCPResourceContent `json:"contents"`
}

// Tool structures
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type MCPListToolsResponse struct {
	Tools []MCPTool `json:"tools"`
}

type MCPCallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	Data string `json:"data,omitempty"`
}

type MCPCallToolResponse struct {
	Result MCPToolResult `json:"result"`
}

// Handle MCP requests
func (s *MCPServer) HandleRequest(c *core.WebContext) (*JSONRPCResponse, *errs.Error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
	}

	// Get current user
	user, err := s.users.GetUserById(c, c.GetCurrentUid())
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	switch req.Method {
	case "initialize":
		return s.handleInitialize(c, &req, user)
	case "resources/list":
		return s.handleListResources(c, &req, user)
	case "resources/read":
		return s.handleReadResource(c, &req, user)
	case "tools/list":
		return s.handleListTools(c, &req, user)
	case "tools/call":
		return s.handleCallTool(c, &req, user)
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
			ID: req.ID,
		}, nil
	}
}

func (s *MCPServer) handleInitialize(c *core.WebContext, req *JSONRPCRequest, user *models.User) (*JSONRPCResponse, *errs.Error) {
	var initReq MCPInitializeRequest
	if req.Params != nil {
		params, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(params, &initReq); err != nil {
			return &JSONRPCResponse{
				JSONRPC: "2.0",
				Error: &JSONRPCError{
					Code:    -32602,
					Message: "Invalid params",
				},
				ID: req.ID,
			}, nil
		}
	}

	response := MCPInitializeResponse{
		ProtocolVersion: MCPVersion,
		Capabilities: MCPCapabilities{
			Resources: &MCPResourceCapabilities{
				Subscribe:   false,
				ListChanged: false,
			},
			Tools: &MCPToolCapabilities{
				ListChanged: false,
			},
		},
		ServerInfo: MCPServerInfo{
			Name:        "ezbookkeeping-mcp",
			Version:     "1.0.0",
			Description: "MCP server for ezBookkeeping transaction data",
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  response,
		ID:      req.ID,
	}, nil
}

func (s *MCPServer) handleListResources(c *core.WebContext, req *JSONRPCRequest, user *models.User) (*JSONRPCResponse, *errs.Error) {
	resources := []MCPResource{
		{
			URI:         "transactions://recent",
			Name:        "Recent Transactions",
			Description: "List of recent transactions",
			MimeType:    "application/json",
		},
		{
			URI:         "accounts://list",
			Name:        "Account List",
			Description: "List of all accounts",
			MimeType:    "application/json",
		},
		{
			URI:         "categories://list",
			Name:        "Category List",
			Description: "List of all transaction categories",
			MimeType:    "application/json",
		},
		{
			URI:         "tags://list",
			Name:        "Tag List",
			Description: "List of all transaction tags",
			MimeType:    "application/json",
		},
	}

	response := MCPListResourcesResponse{
		Resources: resources,
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  response,
		ID:      req.ID,
	}, nil
}

func (s *MCPServer) handleReadResource(c *core.WebContext, req *JSONRPCRequest, user *models.User) (*JSONRPCResponse, *errs.Error) {
	var readReq MCPReadResourceRequest
	if req.Params != nil {
		params, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(params, &readReq); err != nil {
			return &JSONRPCResponse{
				JSONRPC: "2.0",
				Error: &JSONRPCError{
					Code:    -32602,
					Message: "Invalid params",
				},
				ID: req.ID,
			}, nil
		}
	}

	var content string
	var err error

	switch readReq.URI {
	case "transactions://recent":
		content, err = s.getRecentTransactions(c, user)
	case "accounts://list":
		content, err = s.getAccountList(c, user)
	case "categories://list":
		content, err = s.getCategoryList(c, user)
	case "tags://list":
		content, err = s.getTagList(c, user)
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Invalid resource URI",
			},
			ID: req.ID,
		}, nil
	}

	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCError{
				Code:    -32603,
				Message: "Internal error",
				Data:    err.Error(),
			},
			ID: req.ID,
		}, nil
	}

	response := MCPReadResourceResponse{
		Contents: []MCPResourceContent{
			{
				URI:      readReq.URI,
				MimeType: "application/json",
				Text:     content,
			},
		},
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  response,
		ID:      req.ID,
	}, nil
}

func (s *MCPServer) handleListTools(c *core.WebContext, req *JSONRPCRequest, user *models.User) (*JSONRPCResponse, *errs.Error) {
	tools := []MCPTool{
		{
			Name:        "query_transactions",
			Description: "Query transactions with filters",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"start_date": map[string]interface{}{
						"type":        "string",
						"description": "Start date (YYYY-MM-DD)",
					},
					"end_date": map[string]interface{}{
						"type":        "string",
						"description": "End date (YYYY-MM-DD)",
					},
					"category_id": map[string]interface{}{
						"type":        "integer",
						"description": "Category ID to filter by",
					},
					"account_id": map[string]interface{}{
						"type":        "integer",
						"description": "Account ID to filter by",
					},
					"type": map[string]interface{}{
						"type":        "string",
						"description": "Transaction type (income, expense, transfer)",
					},
					"keyword": map[string]interface{}{
						"type":        "string",
						"description": "Keyword to search in transaction comments",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results (default: 100)",
					},
				},
			},
		},
		{
			Name:        "get_transaction_statistics",
			Description: "Get transaction statistics and summaries",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"start_date": map[string]interface{}{
						"type":        "string",
						"description": "Start date (YYYY-MM-DD)",
					},
					"end_date": map[string]interface{}{
						"type":        "string",
						"description": "End date (YYYY-MM-DD)",
					},
					"category_id": map[string]interface{}{
						"type":        "integer",
						"description": "Category ID to filter by",
					},
					"account_id": map[string]interface{}{
						"type":        "integer",
						"description": "Account ID to filter by",
					},
				},
			},
		},
		{
			Name:        "get_account_balance",
			Description: "Get current balance for accounts",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"account_id": map[string]interface{}{
						"type":        "integer",
						"description": "Account ID (optional, returns all accounts if not specified)",
					},
				},
			},
		},
	}

	response := MCPListToolsResponse{
		Tools: tools,
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  response,
		ID:      req.ID,
	}, nil
}

func (s *MCPServer) handleCallTool(c *core.WebContext, req *JSONRPCRequest, user *models.User) (*JSONRPCResponse, *errs.Error) {
	var callReq MCPCallToolRequest
	if req.Params != nil {
		params, _ := json.Marshal(req.Params)
		if err := json.Unmarshal(params, &callReq); err != nil {
			return &JSONRPCResponse{
				JSONRPC: "2.0",
				Error: &JSONRPCError{
					Code:    -32602,
					Message: "Invalid params",
				},
				ID: req.ID,
			}, nil
		}
	}

	var result MCPToolResult
	var err error

	switch callReq.Name {
	case "query_transactions":
		result, err = s.queryTransactions(c, user, callReq.Arguments)
	case "get_transaction_statistics":
		result, err = s.getTransactionStatistics(c, user, callReq.Arguments)
	case "get_account_balance":
		result, err = s.getAccountBalance(c, user, callReq.Arguments)
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCError{
				Code:    -32602,
				Message: "Unknown tool",
			},
			ID: req.ID,
		}, nil
	}

	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCError{
				Code:    -32603,
				Message: "Internal error",
				Data:    err.Error(),
			},
			ID: req.ID,
		}, nil
	}

	response := MCPCallToolResponse{
		Result: result,
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  response,
		ID:      req.ID,
	}, nil
}

// Helper methods for getting data
func (s *MCPServer) getRecentTransactions(c *core.WebContext, user *models.User) (string, error) {
	transactions, err := s.transactions.GetAllTransactions(c, user.Uid, 100, false)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *MCPServer) getAccountList(c *core.WebContext, user *models.User) (string, error) {
	accounts, err := s.accounts.GetAllAccountsByUid(c, user.Uid)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *MCPServer) getCategoryList(c *core.WebContext, user *models.User) (string, error) {
	categories, err := s.transactionCategories.GetAllCategoriesByUid(c, user.Uid, 0, 0)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *MCPServer) getTagList(c *core.WebContext, user *models.User) (string, error) {
	tags, err := s.transactionTags.GetAllTagsByUid(c, user.Uid)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Tool implementations
func (s *MCPServer) queryTransactions(c *core.WebContext, user *models.User, args map[string]interface{}) (MCPToolResult, error) {
	// Parse arguments
	var startTime, endTime int64
	var categoryId, accountId int64
	var transactionType string
	var keyword string
	var limit int32 = 100

	if val, ok := args["start_date"]; ok {
		if startDateStr, ok := val.(string); ok {
			if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
				startTime = utils.GetMinTransactionTimeFromUnixTime(t.Unix())
			}
		}
	}

	if val, ok := args["end_date"]; ok {
		if endDateStr, ok := val.(string); ok {
			if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
				endTime = utils.GetMaxTransactionTimeFromUnixTime(t.Unix())
			}
		}
	}

	if val, ok := args["category_id"]; ok {
		if categoryIdFloat, ok := val.(float64); ok {
			categoryId = int64(categoryIdFloat)
		}
	}

	if val, ok := args["account_id"]; ok {
		if accountIdFloat, ok := val.(float64); ok {
			accountId = int64(accountIdFloat)
		}
	}

	if val, ok := args["type"]; ok {
		if typeStr, ok := val.(string); ok {
			transactionType = typeStr
		}
	}

	if val, ok := args["keyword"]; ok {
		if keywordStr, ok := val.(string); ok {
			keyword = keywordStr
		}
	}

	if val, ok := args["limit"]; ok {
		if limitFloat, ok := val.(float64); ok {
			limit = int32(limitFloat)
		}
	}

	// Set default end time if not provided
	if endTime <= 0 {
		endTime = utils.GetMaxTransactionTimeFromUnixTime(time.Now().Unix())
	}

	// Convert transaction type
	var apiType models.TransactionType
	switch transactionType {
	case "income":
		apiType = models.TRANSACTION_TYPE_INCOME
	case "expense":
		apiType = models.TRANSACTION_TYPE_EXPENSE
	case "transfer":
		apiType = models.TRANSACTION_TYPE_TRANSFER
	default:
		apiType = models.TRANSACTION_TYPE_MODIFY_BALANCE // Use this as "all" types
	}

	// Prepare filter arrays
	var categoryIds []int64
	if categoryId > 0 {
		categoryIds = append(categoryIds, categoryId)
	}

	var accountIds []int64
	if accountId > 0 {
		accountIds = append(accountIds, accountId)
	}

	// Query transactions
	transactions, err := s.transactions.GetAllSpecifiedTransactions(c, user.Uid, endTime, startTime, apiType, categoryIds, accountIds, nil, false, models.TRANSACTION_TAG_FILTER_HAS_ANY, "", keyword, limit, true)
	if err != nil {
		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error querying transactions: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Format result
	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error formatting transactions: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return MCPToolResult{
		Content: []MCPContent{
			{
				Type: "text",
				Text: string(data),
			},
		},
	}, nil
}

func (s *MCPServer) getTransactionStatistics(c *core.WebContext, user *models.User, args map[string]interface{}) (MCPToolResult, error) {
	// Parse arguments
	var startTime, endTime int64
	var categoryId, accountId int64

	if val, ok := args["start_date"]; ok {
		if startDateStr, ok := val.(string); ok {
			if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
				startTime = utils.GetMinTransactionTimeFromUnixTime(t.Unix())
			}
		}
	}

	if val, ok := args["end_date"]; ok {
		if endDateStr, ok := val.(string); ok {
			if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
				endTime = utils.GetMaxTransactionTimeFromUnixTime(t.Unix())
			}
		}
	}

	if val, ok := args["category_id"]; ok {
		if categoryIdFloat, ok := val.(float64); ok {
			categoryId = int64(categoryIdFloat)
		}
	}

	if val, ok := args["account_id"]; ok {
		if accountIdFloat, ok := val.(float64); ok {
			accountId = int64(accountIdFloat)
		}
	}

	// Set default end time if not provided
	if endTime <= 0 {
		endTime = utils.GetMaxTransactionTimeFromUnixTime(time.Now().Unix())
	}

	// Prepare filter arrays
	var categoryIds []int64
	if categoryId > 0 {
		categoryIds = append(categoryIds, categoryId)
	}

	var accountIds []int64
	if accountId > 0 {
		accountIds = append(accountIds, accountId)
	}

	// Get transaction counts
	totalCount, err := s.transactions.GetTransactionCount(c, user.Uid, endTime, startTime, models.TRANSACTION_TYPE_MODIFY_BALANCE, categoryIds, accountIds, nil, false, models.TRANSACTION_TAG_FILTER_HAS_ANY, "", "")
	if err != nil {
		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error getting transaction count: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	incomeCount, err := s.transactions.GetTransactionCount(c, user.Uid, endTime, startTime, models.TRANSACTION_TYPE_INCOME, categoryIds, accountIds, nil, false, models.TRANSACTION_TAG_FILTER_HAS_ANY, "", "")
	if err != nil {
		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error getting income transaction count: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	expenseCount, err := s.transactions.GetTransactionCount(c, user.Uid, endTime, startTime, models.TRANSACTION_TYPE_EXPENSE, categoryIds, accountIds, nil, false, models.TRANSACTION_TAG_FILTER_HAS_ANY, "", "")
	if err != nil {
		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error getting expense transaction count: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	transferCount, err := s.transactions.GetTransactionCount(c, user.Uid, endTime, startTime, models.TRANSACTION_TYPE_TRANSFER, categoryIds, accountIds, nil, false, models.TRANSACTION_TAG_FILTER_HAS_ANY, "", "")
	if err != nil {
		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error getting transfer transaction count: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Create statistics result
	statistics := map[string]interface{}{
		"total_transactions":    totalCount,
		"income_transactions":   incomeCount,
		"expense_transactions":  expenseCount,
		"transfer_transactions": transferCount,
		"period": map[string]interface{}{
			"start_time": startTime,
			"end_time":   endTime,
		},
	}

	// Format result
	data, err := json.MarshalIndent(statistics, "", "  ")
	if err != nil {
		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error formatting statistics: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return MCPToolResult{
		Content: []MCPContent{
			{
				Type: "text",
				Text: string(data),
			},
		},
	}, nil
}

func (s *MCPServer) getAccountBalance(c *core.WebContext, user *models.User, args map[string]interface{}) (MCPToolResult, error) {
	var accountId int64

	if val, ok := args["account_id"]; ok {
		if accountIdFloat, ok := val.(float64); ok {
			accountId = int64(accountIdFloat)
		}
	}

	if accountId > 0 {
		// Get specific account from all accounts
		accounts, err := s.accounts.GetAllAccountsByUid(c, user.Uid)
		if err != nil {
			return MCPToolResult{
				Content: []MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Error getting accounts: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		// Find the specific account
		var targetAccount *models.Account
		for _, account := range accounts {
			if account.AccountId == accountId {
				targetAccount = account
				break
			}
		}

		if targetAccount == nil {
			return MCPToolResult{
				Content: []MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Account with ID %d not found", accountId),
					},
				},
				IsError: true,
			}, nil
		}

		data, err := json.MarshalIndent(targetAccount, "", "  ")
		if err != nil {
			return MCPToolResult{
				Content: []MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Error formatting account: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: string(data),
				},
			},
		}, nil
	} else {
		// Get all accounts
		accounts, err := s.accounts.GetAllAccountsByUid(c, user.Uid)
		if err != nil {
			return MCPToolResult{
				Content: []MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Error getting accounts: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		data, err := json.MarshalIndent(accounts, "", "  ")
		if err != nil {
			return MCPToolResult{
				Content: []MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Error formatting accounts: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		return MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: string(data),
				},
			},
		}, nil
	}
}