package mcp

import (
	"reflect"

	"github.com/invopop/jsonschema"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/services"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// MCPAvailableServices holds the services available for MCP tools
type MCPAvailableServices interface {
	GetTransactionService() *services.TransactionService
	GetTransactionCategoryService() *services.TransactionCategoryService
	GetTransactionTagService() *services.TransactionTagService
	GetAccountService() *services.AccountService
	GetUserService() *services.UserService
}

// MCPToolHandler defines the MCP tool handler
type MCPToolHandler[T MCPTextContent | MCPImageContent | MCPAudioContent | MCPResourceLink | MCPEmbeddedResource] interface {
	// Description returns the description of the MCP tool
	Description() string

	// InputType returns the input type for the MCP tool request
	InputType() reflect.Type

	// OutputType returns the output type for the MCP tool response
	OutputType() reflect.Type

	// Handle processes the MCP call tool request and returns the response
	Handle(*core.WebContext, *MCPCallToolRequest, *settings.Config, MCPAvailableServices) ([]*T, *errs.Error)
}

// GetAllMCPToolInfos returns all available MCP tool information
func GetAllMCPToolInfos() []*MCPTool {
	toolInfos := make([]*MCPTool, 0)

	for name, handler := range mcpTextContentTools {
		toolInfos = append(toolInfos, getMCPToolInfo(name, handler))
	}

	for name, handler := range mcpImageContentTools {
		toolInfos = append(toolInfos, getMCPToolInfo(name, handler))
	}

	for name, handler := range mcpAudioContentTools {
		toolInfos = append(toolInfos, getMCPToolInfo(name, handler))
	}

	for name, handler := range mcpResourceLinkTools {
		toolInfos = append(toolInfos, getMCPToolInfo(name, handler))
	}

	for name, handler := range mcpEmbeddedResourceTools {
		toolInfos = append(toolInfos, getMCPToolInfo(name, handler))
	}

	return toolInfos
}

// MCPToolHandle handles the MCP tool request based on the tool name
func MCPToolHandle(c *core.WebContext, callToolReq *MCPCallToolRequest, currentConfig *settings.Config, services MCPAvailableServices) (any, *errs.Error) {
	if handler, exists := mcpTextContentTools[callToolReq.Name]; exists {
		return mcpTextContentToolHandle(c, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := mcpImageContentTools[callToolReq.Name]; exists {
		return mcpImageContentToolHandle(c, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := mcpAudioContentTools[callToolReq.Name]; exists {
		return mcpAudioContentToolHandle(c, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := mcpResourceLinkTools[callToolReq.Name]; exists {
		return mcpResourceLinkToolHandle(c, handler, currentConfig, services, callToolReq)
	}

	if handler, exists := mcpEmbeddedResourceTools[callToolReq.Name]; exists {
		return mcpEmbeddedResourceToolHandle(c, handler, currentConfig, services, callToolReq)
	}

	return nil, errs.ErrApiNotFound
}

func mcpTextContentToolHandle(c *core.WebContext, handler MCPToolHandler[MCPTextContent], currentConfig *settings.Config, services MCPAvailableServices, callToolReq *MCPCallToolRequest) (any, *errs.Error) {
	result, err := handler.Handle(c, callToolReq, currentConfig, services)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	callToolResp := MCPCallToolResponse[MCPTextContent]{
		Content: result,
		IsError: false,
	}

	return callToolResp, nil
}

func mcpImageContentToolHandle(c *core.WebContext, handler MCPToolHandler[MCPImageContent], currentConfig *settings.Config, services MCPAvailableServices, callToolReq *MCPCallToolRequest) (any, *errs.Error) {
	result, err := handler.Handle(c, callToolReq, currentConfig, services)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	callToolResp := MCPCallToolResponse[MCPImageContent]{
		Content: result,
		IsError: false,
	}

	return callToolResp, nil
}

func mcpAudioContentToolHandle(c *core.WebContext, handler MCPToolHandler[MCPAudioContent], currentConfig *settings.Config, services MCPAvailableServices, callToolReq *MCPCallToolRequest) (any, *errs.Error) {
	result, err := handler.Handle(c, callToolReq, currentConfig, services)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	callToolResp := MCPCallToolResponse[MCPAudioContent]{
		Content: result,
		IsError: false,
	}

	return callToolResp, nil
}

func mcpResourceLinkToolHandle(c *core.WebContext, handler MCPToolHandler[MCPResourceLink], currentConfig *settings.Config, services MCPAvailableServices, callToolReq *MCPCallToolRequest) (any, *errs.Error) {
	result, err := handler.Handle(c, callToolReq, currentConfig, services)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	callToolResp := MCPCallToolResponse[MCPResourceLink]{
		Content: result,
		IsError: false,
	}

	return callToolResp, nil
}

func mcpEmbeddedResourceToolHandle(c *core.WebContext, handler MCPToolHandler[MCPEmbeddedResource], currentConfig *settings.Config, services MCPAvailableServices, callToolReq *MCPCallToolRequest) (any, *errs.Error) {
	result, err := handler.Handle(c, callToolReq, currentConfig, services)

	if err != nil {
		return nil, errs.Or(err, errs.ErrOperationFailed)
	}

	callToolResp := MCPCallToolResponse[MCPEmbeddedResource]{
		Content: result,
		IsError: false,
	}

	return callToolResp, nil
}

func getMCPToolInfo[T MCPTextContent | MCPImageContent | MCPAudioContent | MCPResourceLink | MCPEmbeddedResource](name string, handler MCPToolHandler[T]) *MCPTool {
	mcpTool := &MCPTool{
		Name:        name,
		Description: handler.Description(),
	}

	schemeGenerator := jsonschema.Reflector{
		Anonymous:      true,
		DoNotReference: true,
		ExpandedStruct: true,
	}

	if handler.InputType() != nil {
		schema := schemeGenerator.ReflectFromType(handler.InputType())
		mcpTool.InputSchema = schema
	}

	if handler.OutputType() != nil {
		schema := schemeGenerator.ReflectFromType(handler.OutputType())
		mcpTool.OutputSchema = schema
	}

	return mcpTool
}
