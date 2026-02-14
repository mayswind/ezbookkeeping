package api

import (
	"encoding/json"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/mcp"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const mcpServerName = core.ApplicationName + "-mcp"

// ModelContextProtocolAPI represents model context protocol api
type ModelContextProtocolAPI struct {
	ApiUsingConfig
	transactions          *services.TransactionService
	transactionCategories *services.TransactionCategoryService
	transactionTags       *services.TransactionTagService
	accounts              *services.AccountService
	users                 *services.UserService
	tokens                *services.TokenService
}

// Initialize a model context protocol api singleton instance
var (
	ModelContextProtocols = &ModelContextProtocolAPI{
		ApiUsingConfig: ApiUsingConfig{
			container: settings.Container,
		},
		transactions:          services.Transactions,
		transactionCategories: services.TransactionCategories,
		transactionTags:       services.TransactionTags,
		accounts:              services.Accounts,
		users:                 services.Users,
		tokens:                services.Tokens,
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

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Warnf(c, "[model_context_protocols.InitializeHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_MCP_ACCESS) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	tokenClaims := c.GetTokenClaims()
	userTokenId, err := utils.StringToInt64(tokenClaims.UserTokenId)

	if err != nil {
		log.Warnf(c, "[model_context_protocols.InitializeHandler] parse user token id failed, because %s", err.Error())
	} else {
		tokenRecord := &models.TokenRecord{
			Uid:             tokenClaims.Uid,
			UserTokenId:     userTokenId,
			CreatedUnixTime: tokenClaims.IssuedAt,
		}

		tokenId := a.tokens.GenerateTokenId(tokenRecord)

		err = a.tokens.UpdateTokenLastSeen(c, tokenRecord)

		if err != nil {
			log.Warnf(c, "[model_context_protocols.InitializeHandler] failed to update last seen of token \"id:%s\" for user \"uid:%d\", because %s", tokenId, uid, err.Error())
		}
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
			Title:   core.ApplicationName,
			Version: core.Version,
		},
	}

	return initResp, nil
}

// ListResourcesHandler returns the list of resources for model context protocol
func (a *ModelContextProtocolAPI) ListResourcesHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Warnf(c, "[model_context_protocols.ListResourcesHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_MCP_ACCESS) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

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

	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Warnf(c, "[model_context_protocols.ReadResourceHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_MCP_ACCESS) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	return nil, errs.ErrApiNotFound
}

// ListToolsHandler returns the list of tools for model context protocol
func (a *ModelContextProtocolAPI) ListToolsHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Warnf(c, "[model_context_protocols.ListToolsHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_MCP_ACCESS) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	mcpVersion := a.getMCPVersion(c)
	toolsInfo := mcp.Container.GetMCPTools()
	finalToolsInfos := make([]*mcp.MCPTool, len(toolsInfo))

	for i := 0; i < len(toolsInfo); i++ {
		finalToolsInfos[i] = &mcp.MCPTool{
			Name:        toolsInfo[i].Name,
			InputSchema: toolsInfo[i].InputSchema,
			Title:       toolsInfo[i].Title,
			Description: toolsInfo[i].Description,
		}

		if mcpVersion >= string(mcp.ToolResultStructuredContentMinVersion) {
			finalToolsInfos[i].OutputSchema = toolsInfo[i].OutputSchema
		}
	}

	listToolsResp := mcp.MCPListToolsResponse{
		Tools: finalToolsInfos,
	}

	return listToolsResp, nil
}

// CallToolHandler returns the result of calling a specific tool for model context protocol
func (a *ModelContextProtocolAPI) CallToolHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	uid := c.GetCurrentUid()
	user, err := a.users.GetUserById(c, uid)

	if err != nil {
		log.Warnf(c, "[model_context_protocols.CallToolHandler] failed to get user \"uid:%d\" info, because %s", uid, err.Error())
		return nil, errs.ErrUserNotFound
	}

	if user.FeatureRestriction.Contains(core.USER_FEATURE_RESTRICTION_TYPE_MCP_ACCESS) {
		return nil, errs.ErrNotPermittedToPerformThisAction
	}

	var callToolReq mcp.MCPCallToolRequest

	if jsonRPCRequest.Params != nil {
		if err := json.Unmarshal(jsonRPCRequest.Params, &callToolReq); err != nil {
			return nil, errs.NewIncompleteOrIncorrectSubmissionError(err)
		}
	} else {
		return nil, errs.ErrIncompleteOrIncorrectSubmission
	}

	result, err := mcp.Container.HandleTool(c, &callToolReq, user, a.CurrentConfig(), a)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	return result, nil
}

// PingHandler return the ping response for model context protocol
func (a *ModelContextProtocolAPI) PingHandler(c *core.WebContext, jsonRPCRequest *core.JSONRPCRequest) (any, *errs.Error) {
	return core.O{}, nil
}

// GetTransactionService implements the MCPAvailableServices interface
func (a *ModelContextProtocolAPI) GetTransactionService() *services.TransactionService {
	return a.transactions
}

// GetTransactionCategoryService implements the MCPAvailableServices interface
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

// GetUserService implements the MCPAvailableServices interface
func (a *ModelContextProtocolAPI) GetUserService() *services.UserService {
	return a.users
}

// getMCPVersion returns the MCP protocol version from the request header
func (a *ModelContextProtocolAPI) getMCPVersion(c *core.WebContext) string {
	return c.GetHeader(mcp.MCPProtocolVersionHeaderName)
}
