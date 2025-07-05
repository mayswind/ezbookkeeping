package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/mcp"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

const mcpServerName = "ezBookkeeping-mcp"

// ModelContextProtocolAPI represents model context protocol api
type ModelContextProtocolAPI struct {
	ApiUsingConfig
	transactions            *services.TransactionService
	transactionCategories   *services.TransactionCategoryService
	transactionTags         *services.TransactionTagService
	accounts                *services.AccountService
	users                   *services.UserService
	userCustomExchangeRates *services.UserCustomExchangeRatesService
}

// Initialize a model context protocol api singleton instance
var (
	ModelContextProtocols = &ModelContextProtocolAPI{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		transactions:            services.Transactions,
		transactionCategories:   services.TransactionCategories,
		transactionTags:         services.TransactionTags,
		accounts:                services.Accounts,
		users:                   services.Users,
		userCustomExchangeRates: services.UserCustomExchangeRates,
	}
)

// InitializeHandler returns the initialize response for model context protocol
func (a *ModelContextProtocolAPI) InitializeHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	var initRequest mcp.MCPInitializeRequest

	if jsonRPCRequest.Params != nil {
		if err := json.Unmarshal(jsonRPCRequest.Params, &initRequest); err != nil {
			return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}
	} else {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	protocolVersion := mcp.MCPProtocolVersion(initRequest.ProtocolVersion)
	_, exists := mcp.SupportedMCPVersion[protocolVersion]

	if !exists {
		protocolVersion = mcp.LatestSupportedMCPVersion
	}

	initResp := mcp.MCPInitializeResponse{
		ProtocolVersion: string(protocolVersion),
		Capabilities: &mcp.MCPCapabilities{
			Tools: &mcp.MCPToolCapabilities{
				ListChanged: false,
			},
		},
		ServerInfo: &mcp.MCPImplementation{
			Name:    mcpServerName,
			Title:   a.CurrentConfig().AppName,
			Version: settings.Version,
		},
	}

	return initResp, nil
}

// ListResourcesHandler returns the list of resources for model context protocol
func (a *ModelContextProtocolAPI) ListResourcesHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	listResourcesResp := mcp.MCPListResourcesResponse{
		Resources: make([]*mcp.MCPResource, 0),
	}

	return listResourcesResp, nil
}

// ReadResourceHandler returns the resource details for a specific resource in model context protocol
func (a *ModelContextProtocolAPI) ReadResourceHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	var readResourceReq mcp.MCPReadResourceRequest

	if jsonRPCRequest.Params != nil {
		if err := json.Unmarshal(jsonRPCRequest.Params, &readResourceReq); err != nil {
			return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}
	} else {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	return nil, errs.ErrApiNotFound
}

// ListToolsHandler returns the list of tools for model context protocol
func (a *ModelContextProtocolAPI) ListToolsHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	listToolsResp := mcp.MCPListToolsResponse{
		Tools: mcp.AllMCPToolInfos,
	}

	return listToolsResp, nil
}

// CallToolHandler returns the result of calling a specific tool for model context protocol
func (a *ModelContextProtocolAPI) CallToolHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	var callToolReq mcp.MCPCallToolRequest

	if jsonRPCRequest.Params != nil {
		if err := json.Unmarshal(jsonRPCRequest.Params, &callToolReq); err != nil {
			return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}
	} else {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	result, err := mcp.MCPToolHandle(c, &callToolReq, a.CurrentConfig(), a)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// PingHandler return the ping response for model context protocol
func (a *ModelContextProtocolAPI) PingHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	return gin.H{}, nil
}

// GetTransactionService implements the MCPAvailableServices interface
func (a *ModelContextProtocolAPI) GetTransactionService() *services.TransactionService {
	return a.transactions
}

// GetUserCustomExchangeRatesService implements the MCPAvailableServices interface
func (a *ModelContextProtocolAPI) GetTransactionCategoryService() *services.TransactionCategoryService {
	return a.transactionCategories
}

// GetTransactionTagService implements the MCPAvailableServices interface
func (a *ModelContextProtocolAPI) GetTransactionTagService() *services.TransactionTagService {
	return a.transactionTags
}

// GetAccountService implements the MCPAvailableServices interface
func (a *ModelContextProtocolAPI) GetAccountService() *services.AccountService {
	return a.accounts
}

// GetUserCustomExchangeRatesService implements the MCPAvailableServices interface
func (a *ModelContextProtocolAPI) GetUserService() *services.UserService {
	return a.users
}
